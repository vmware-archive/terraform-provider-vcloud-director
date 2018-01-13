package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/grpc"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	"os"
	"runtime/debug"
)

type CatalogItem interface {
	create() error
}

type CommonResource struct {
	d *schema.ResourceData
	p grpc.PyVcloudProvider
}

type CatalogItemBase struct {
	CatalogName string
	ItemName    string

	//common reference
	resource CommonResource
}
type CatalogItemOvaFromFile struct {
	catalogItem CatalogItemBase

	FilePath string
}

type CatalogItemOvaFromVApp struct {
	catalogItem CatalogItemBase

	VappName               string
	VdcName                string
	Description            string
	CustomizeOnInstantiate bool
}

func (c CatalogItemOvaFromFile) create() error {
	logging.Plogf("__INIT__create for %+v", c)
	catalogName, itemName, filePath := c.catalogItem.CatalogName, c.catalogItem.ItemName, c.FilePath
	d := c.catalogItem.resource.d
	provider := c.catalogItem.resource.p

	logging.Plogf("__LOG__checking for %v", filePath)
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {

		debug.PrintStack()
		//time.Sleep(20000 * time.Millisecond)
		os.Exit(1001)

		return fmt.Errorf("__ERROR__ File does not exist FILE =  [%v] ", filePath)

	}

	logging.Plogf("__LOG__File Present_ Path = [%v]", filePath)

	catalogIsPresInfo := proto.IsPresentCatalogItemInfo{CatalogName: catalogName, ItemName: itemName}

	isPresResp, isPreErr := provider.IsPresentCatalogItem(catalogIsPresInfo)

	if isPreErr != nil {
		return fmt.Errorf("Error Creating Item :[%v] %#v", catalogName, isPreErr)
	}
	if isPresResp.Present {
		logging.PlogWarn(fmt.Sprintf("__LOG__ catalog item [%v] is already present /setting state as created ", catalogName))
		d.SetId(getOvaId(d))
		return nil
	}

	catalogUploadInfo := proto.CatalogUploadOvaInfo{CatalogName: catalogName,
		FilePath: filePath,
		ItemName: itemName,
	}

	resp, errp := provider.CatalogUploadOva(catalogUploadInfo)
	if errp != nil {
		logging.PlogErrorf("__ERROR__ creating catalog ITEM %v", errp)
		return fmt.Errorf("__ERROR__ Creating catalog Item: [%#v]", errp)
	}

	logging.Plog(fmt.Sprintf("__LOG__ resp.Created [%v]", resp.Created))
	//WAIT FOR STATUS to RESOLVED
	// https://github.com/vmware/container-service-extension/blob/1a4d2d38a2243dcc4cf2eac1cdd4d75335122e72/container_service_extension/config.py#L348

	checkResolved := proto.CatalogCheckResolvedInfo{
		CatalogName: catalogName,
		ItemName:    itemName,
	}

	cres, errc := provider.OvaCheckResolved(checkResolved)
	if errc != nil {
		logging.Plog("__ERROR__ OvaCheckResolved catalog ITEM ")
		return fmt.Errorf("Error OvaCheckResolved catalog Item: [%#v]", errc)
	}
	if resp.Created && cres.Resolved {
		logging.Plog(fmt.Sprintf("__LOG__catalog Item [%v] is Created and RESOLVED  ", getOvaId(d)))
		d.SetId(getOvaId(d))
		return nil
	}
	return nil

}

//is present
func (c CatalogItemBase) isPresent() (bool, error) {
	logging.Plogf("__INIT__isPresent for %+v", c)
	catalogIsPresInfo := proto.IsPresentCatalogItemInfo{CatalogName: c.CatalogName, ItemName: c.ItemName}
	provider, d := c.resource.p, c.resource.d

	isPresResp, isPreErr := provider.IsPresentCatalogItem(catalogIsPresInfo)

	if isPreErr != nil {
		return false, fmt.Errorf("Error Creating Item :[%v] %#v", c.CatalogName, isPreErr)
	}
	if isPresResp.Present {
		logging.PlogWarn(fmt.Sprintf("__LOG__ catalog item [%v] is already present /setting state as created ", c.CatalogName))
		d.SetId(getOvaId(d))
		return true, nil

	}
	logging.Plog(fmt.Sprintf("__DONE__isPresent for [%+v] ===> [%+v]  ", c, isPresResp))
	return false, nil

}

//capture vapp
func (c CatalogItemOvaFromVApp) create() error {
	logging.Plogf("__INIT__create for %+v", c)
	catalogName, itemName := c.catalogItem.CatalogName, c.catalogItem.ItemName
	provider, d := c.catalogItem.resource.p, c.catalogItem.resource.d

	present, err := c.catalogItem.isPresent()
	if err != nil {
		return err
	}

	logging.Plogf("present - %v", present)
	if present {
		//is present has handled return the flow
		logging.Plogf("__DONE EARLY EXIT __create for %+v", c)
		return nil
	}

	captureInfo := proto.CaptureVAppInfo{
		CatalogName:            catalogName,
		ItemName:               itemName,
		VappName:               c.VappName,
		VdcName:                c.VdcName,
		CustomizeOnInstantiate: c.CustomizeOnInstantiate,
	}

	resp, errp := provider.CaptureVapp(captureInfo)
	logging.Plog(fmt.Sprintf("__LOG__ CaptureVapp completed %#v  ", errp))

	if errp != nil {
		logging.PlogErrorf("__ERROR__ Creating catalog Item: [%#v]", errp)
		return fmt.Errorf("__ERROR__ Creating catalog Item: [%#v]", errp)
	}

	checkResolved := proto.CatalogCheckResolvedInfo{
		CatalogName: catalogName,
		ItemName:    itemName,
	}

	cres, errc := provider.OvaCheckResolved(checkResolved)
	if errc != nil {
		logging.Plog("__ERROR__ OvaCheckResolved catalog ITEM ")
		return fmt.Errorf("Error OvaCheckResolved catalog Item: [%#v]", errc)
	}
	if resp.Captured && cres.Resolved {
		logging.Plog(fmt.Sprintf("__LOG__catalog Item [%v] is Created and RESOLVED  ", getOvaId(d)))
		d.SetId(getOvaId(d))
		return nil
	}
	logging.Plogf("__DONE__create for %+v", c)
	return nil
}

func buildCatalogItem(d *schema.ResourceData, m interface{}) CatalogItem {

	itemName := d.Get("item_name").(string)
	catalogName := d.Get("catalog_name").(string)
	source_file_path := d.Get("source_file_path").(string)
	source_vapp_name := d.Get("source_vapp_name").(string)
	source_vdc_name := d.Get("source_vdc_name").(string)
	customize_on_instantiate := d.Get("customize_on_instantiate").(bool)
	logging.Plog(fmt.Sprintf("buildCatalogItem [%v] [%v] [%v] [%v] [%v]", itemName, catalogName, source_file_path, source_vapp_name, source_vdc_name))
	r := CommonResource{d: d, p: getProvider(m)}

	cb := CatalogItemBase{CatalogName: catalogName, ItemName: itemName, resource: r}
	if len(source_file_path) == 0 {
		logging.Plog("create from VAPP")
		c := CatalogItemOvaFromVApp{
			catalogItem:            cb,
			VappName:               source_vapp_name,
			VdcName:                source_vdc_name,
			CustomizeOnInstantiate: customize_on_instantiate,
		}
		return c

	} else {
		c := CatalogItemOvaFromFile{
			catalogItem: cb,
			FilePath:    source_file_path,
		}

		return c
	}

}
