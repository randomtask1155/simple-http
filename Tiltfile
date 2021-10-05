
k8s_yaml('./workload.yml')
k8s_kind('Workload', image_json_path='{.metadata.annotations.image}')
k8s_resource(workload='simple-http', extra_pod_selectors=[{'service.knatvie.dev/service':'simple-http'}])