#go test disk_test.go -v run  TestDiskInterface

export TF_ACC=1






#go test github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/provider/ -v -run TestAccResourceIndependentDiskBasic 

#| grep --line-buffered -vE 'TRACE|terraform|^$'
go test test -v -run TestDiskInterface 