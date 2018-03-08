import logging

from pyvcloud.vcd.client import BasicLoginCredentials
from pyvcloud.vcd.client import Client
from pyvcloud.vcd.client import TaskStatus
from pyvcloud.vcd.org import Org
from pyvcloud.vcd.vapp import VApp
from pyvcloud.vcd.vdc import VDC
#from errors import VappVmCreationError
from pyvcloud.vcd.vm import VM


def vcd_login(host, user, password, org):
    logging.basicConfig(level=logging.DEBUG)
    logging.info("login called")
    client = Client(
        host,
        api_version="27.0",
        verify_ssl_certs=False,
        log_file='vcd.log',
        log_requests=True,
        log_headers=True,
        log_bodies=True)
    try:
        client.set_credentials(BasicLoginCredentials(user, org, password))

        return client
    except Exception as e:
        print('error occured', e)


def create(client):
    print("=============== __LOG__Create_VDC  =======================\n\n")

    vdc_name = "ACME_PAYG"
    vapp_name = "test2"

    org_resource = client.get_org()
    org = Org(client, resource=org_resource)
    print("Org name: ", org.get_name())
    print("Vdc name: ", vdc_name)

    try:
        vdc_resource = org.get_vdc(vdc_name)
        vdc = VDC(client, name=vdc_name, resource=vdc_resource)

        vapp_resource = vdc.get_vapp(vapp_name)
        vapp = VApp(client, name=vapp_name, resource=vapp_resource)
        print("vapp : ", vapp)

        catalog_item = org.get_catalog_item('ACME', 'tinyova')

        source_vapp_resource = client.get_resource(
            catalog_item.Entity.get('href'))

        print("source_vapp_resource: ", source_vapp_resource)

        spec = {
            'source_vm_name': 'Tiny Linux template',
            'vapp': source_vapp_resource
        }

        storage_profiles = [{
            'name': 'Performance',
            'enabled': True,
            'units': 'MB',
            'limit': 0,
            'default': True
        }]

        spec['target_vm_name'] = 'ubuntu_pcp_11'
        spec['hostname'] = 'ubuntu'
        spec['network'] = 'global'
        spec['ip_allocation_mode'] = 'dhcp'
        #spec['storage_profile'] = storage_profiles

        vms = [spec]
        result = vapp.add_vms(vms)

        print("result: ", result)
        #task = client.get_task_monitor().wait_for_status(
        #                    task=result,
        #                    timeout=60,
        #                    poll_frequency=2,
        #                    fail_on_statuses=None,
        #                    expected_target_statuses=[
        #                        TaskStatus.SUCCESS,
        #                        TaskStatus.ABORTED,
        #                        TaskStatus.ERROR,
        #                        TaskStatus.CANCELED],
        #                    callback=None)

        #st = task.get('status')
        #if st == TaskStatus.SUCCESS.value:
        #    message = 'status : {0} '.format(st)
        #    logging.info(message)
        #else:
        #    print("st : ", st)
        #    raise Exception(task)
        print("=============================================\n\n")
        return True
    except Exception as e:
        error_message = '__ERROR_ [create_vdc] failed for vdc {0} '.format(
            vdc_name)
        logging.warn(error_message, e)
        return False


def read_vdc(client):
    print("================= Vdc read request ===================")
    vdc_name = "pcp_vdc_03"
    org_resource = client.get_org()
    org = Org(client, resource=org_resource)
    print("Org name: ", org.get_name())
    print("Vdc name: ", vdc_name)

    vdc_resource = org.get_vdc(vdc_name)
    vdc = VDC(client, name=vdc_name, resource=vdc_resource)

    print("name = ", vdc_resource.get('name'))
    # res.provider_vdc = str(vdc_resource.provider_vdc)
    description = str(vdc_resource.Description)
    print("description = ", description)

    allocation_model = str(vdc_resource.AllocationModel)
    print("allocation_model = ", allocation_model)

    cpu_units = str(vdc_resource.ComputeCapacity.Cpu.Units)
    print("cpu_units = ", cpu_units)

    cpu_allocated = vdc_resource.ComputeCapacity.Cpu.Allocated
    print("cpu_allocated = ", cpu_allocated)

    cpu_limit = vdc_resource.ComputeCapacity.Cpu.Limit
    print("cpu_limit = ", cpu_limit)

    mem_units = vdc_resource.ComputeCapacity.Memory.Units
    print("mem_units = ", mem_units)

    mem_allocated = vdc_resource.ComputeCapacity.Memory.Allocated
    print("mem_allocated = ", mem_allocated)

    mem_limit = vdc_resource.ComputeCapacity.Memory.Limit
    print("mem_limit = ", mem_limit)

    nic_quota = vdc_resource.NicQuota
    print("nic_quota = ", nic_quota)

    network_quota = vdc_resource.NetworkQuota
    print("network_quota = ", network_quota)

    vm_quota = vdc_resource.VmQuota
    print("vm_quota = ", vm_quota)

    storage_profiles = str(
        vdc_resource.VdcStorageProfiles.VdcStorageProfile.get('name'))
    print("storage_profiles = ", storage_profiles)

    # res.resource_guaranteed_memory = str(vdc_resource.resource_guaranteed_memory)
    # res.resource_guaranteed_cpu = str(vdc_resource.resource_guaranteed_cpu)

    vcpu_in_mhz = vdc_resource.VCpuInMhz2
    print("vcpu_in_mhz = ", vcpu_in_mhz)

    # res.is_thin_provision = str(vdc_resource.is_thin_provision)
    # res.network_pool_name = str(vdc_resource.network_pool_name)
    # res.uses_fast_provisioning = str(vdc_resource.uses_fast_provisioning)
    # res.over_commit_allowed = str(vdc_resource.over_commit_allowed)
    # res.vm_discovery_enabled = str(vdc_resource.vm_discovery_enabled)

    is_enabled = vdc_resource.IsEnabled
    print("is_enabled = ", is_enabled)


def update_vdc(client):
    print("================= Vdc update request ===================")
    vdc_name = "pcp_vdc_02"
    is_enabled = False

    org_resource = client.get_org()
    org = Org(client, resource=org_resource)
    print("Org name: ", org.get_name())
    print("Vdc name: ", vdc_name)

    vdc_resource = org.get_vdc(vdc_name)
    vdc = VDC(client, name=vdc_name, resource=vdc_resource)
    update_vdc_resp = vdc.enable_vdc(is_enabled)

    print("================= Vdc updated ===================")


def delete_vm(client):
    print("================= Vdc delete request ===================")
    vdc_name = "pcp_vdc_02"
    target_vm_name = "pcp_vm"
    org_resource = client.get_org()
    org = Org(client, resource=org_resource)
    print("Org name: ", org.get_name())
    print("Vdc name: ", vdc_name)

    vdc_resource = org.get_vdc(vdc_name)
    vdc = VDC(client, name=vdc_name, resource=vdc_resource)

    vapp_resource = vdc.get_vapp(vapp_name)
    vapp = VApp(client, name=vapp_name, resource=vapp_resource)

    delete_vapp_vm_resp = vapp.delete_vms(target_vm_name)
    task = client.get_task_monitor().wait_for_status(
        task=delete_vapp_vm_resp,
        timeout=60,
        poll_frequency=2,
        fail_on_statuses=None,
        expected_target_statuses=[
            TaskStatus.SUCCESS, TaskStatus.ABORTED, TaskStatus.ERROR,
            TaskStatus.CANCELED
        ],
        callback=None)

    st = task.get('status')
    if st == TaskStatus.SUCCESS.value:
        message = 'delete vdc status : {0} '.format(st)
        logging.info(message)
    else:
        raise errors.VCDVdcDeleteError(etree.tostring(task, pretty_print=True))


def power_off(client, target_vdc, target_vapp, target_vm_name):
    logging.info("__INIT__delete[VappVm]")

    powered_off = False
    org_resource = client.get_org()
    org = Org(client, resource=org_resource)
    try:
        vdc_resource = org.get_vdc(target_vdc)
        vdc = VDC(slef.client, name=target_vdc, resource=vdc_resource)

        vapp_resource = vdc.get_vapp(target_vapp)
        vapp = VApp(client, name=target_vapp, resource=vapp_resource)
        vapp_vm_resource = vapp.get_vm(target_vm_name)
        vm = VM(client, resource=vapp_vm_resource)
        vm.power_off()
        powered_off = True
    except Exception as e:
        error_message = '__ERROR_power_off[VappVm] failed for VappVm {0}. __ErrorMessage__ {1}'.format(
            target_vm_name, str(e))
        logging.warn(error_message)

    return powered_off


if __name__ == '__main__':
    # client = vcd_login("csa-sandbox.eng.vmware.com", "ndixit", "Shephertz@12345", "ext")
    # client = vcd_login("10.172.158.119", "administrator", "VMware1!", "acme")
    client = vcd_login("10.172.158.127", "acmeadmin", "VMware1!", "acme")

    create(client)
    # #delete_vdc(client)
    # #update_vdc(client)
    # read_vdc(client)
    power_off(client, "ACME_PAYG", "test2", "pcp_hi_05")
