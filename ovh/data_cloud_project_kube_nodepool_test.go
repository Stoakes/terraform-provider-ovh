package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudProjectKubeNodePoolDataSource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_KUBE_REGION_TEST")

	config := fmt.Sprintf(
		testAccCloudProjectKubeNodePoolDataSourceConfig,
		os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"),
		name,
		region,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKubernetes(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_kube_nodepool.poolDataSource", "max_nodes", "2"),
				),
			},
		},
	})
}

var testAccCloudProjectKubeNodePoolDataSourceConfig = `
resource "ovh_cloud_project_kube" "cluster" {
	service_name  = "%s"
	name          = "%s"
	region        = "%s"
}

resource "ovh_cloud_project_kube_nodepool" "pool" {
	service_name  = ovh_cloud_project_kube.cluster.service_name
	kube_id       = ovh_cloud_project_kube.cluster.id
	name          = ovh_cloud_project_kube.cluster.name
	flavor_name   = "b2-7"
	desired_nodes = 1
	min_nodes     = 0
	max_nodes     = 2

	depends_on = [
		ovh_cloud_project_kube.cluster
	]
}

data "ovh_cloud_project_kube_nodepool" "poolDataSource" {
  service_name  = ovh_cloud_project_kube.cluster.service_name
  kube_id       = ovh_cloud_project_kube.cluster.id
  name          = ovh_cloud_project_kube_nodepool.pool.name

  depends_on = [
    ovh_cloud_project_kube_nodepool.pool
  ]
}
`
