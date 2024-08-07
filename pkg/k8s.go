package pkg

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	appsV1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	networkV1 "k8s.io/api/networking/v1"
	storageV1 "k8s.io/api/storage/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	rbacV1 "k8s.io/api/rbac/v1"
	//rbacV1beta1 "k8s.io/api/rbac/v1beta1"
)

type K8sClient struct {
	Client *kubernetes.Clientset
}

func NewK8sClient(base64kubeconfig string) (*K8sClient,error){
	decoded, _ := base64.StdEncoding.DecodeString(base64kubeconfig)
	// 使用kubeconfig文件来获取客户端配置
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(decoded))
	if err != nil {
		return nil,err 
	}
 
	// 根据客户端配置创建一个Kubernetes客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil,err 
	}
 
	return &K8sClient{
		Client:clientset,
	},nil 
}

func (k *K8sClient) CreateNs(ns string) error {
	// 创建namespace的spec
	var namespace coreV1.Namespace
	namespace.Name = ns
	defer k.CloseClient()
	_, err := k.Client.CoreV1().Namespaces().Create(context.TODO(),&namespace,metaV1.CreateOptions{})
    if err != nil {
        return err 
    }
	return nil 
}
func (k *K8sClient) UpdateSecretByOpaque(ns,secretname string,kv map[string][]byte) error {
	defer k.CloseClient()
	secret := &coreV1.Secret{
        ObjectMeta: metaV1.ObjectMeta{
            Name: secretname,
        },
        Type: "Opaque",
        Data: kv,
    }
	_,err := k.Client.CoreV1().Secrets(ns).Update(context.TODO(), secret, metaV1.UpdateOptions{})
    if err != nil {
        return err 
    }
	return nil 
}
/*
func (k *K8sClient) ClusterUsers() (error) {
	defer k.CloseClient()
	// 获取用户(ServiceAccount)列表
    users, err := k.Client.RbacV1().Users().List(context.TODO(), metaV1.ListOptions{})
    if err != nil {
       	return err 
    }
	return users,nil 
}
*/
func (k *K8sClient) CreateSecretByOpaque(ns,secretname string,kv map[string][]byte) error {
	defer k.CloseClient()
	secret := &coreV1.Secret{
        ObjectMeta: metaV1.ObjectMeta{
            Name: secretname,
        },
        Type: "Opaque",
        Data: kv,
    }
	_,err := k.Client.CoreV1().Secrets(ns).Create(context.TODO(), secret, metaV1.CreateOptions{})
    if err != nil {
        return err 
    }
	return nil 
}
func (k *K8sClient) UpdateSecretByTlsCert(ns,secretname,cert,key string) error {
	defer k.CloseClient()
	data := map[string][]byte{
        coreV1.TLSCertKey:       []byte(cert),
        coreV1.TLSPrivateKeyKey: []byte(key),
    }
	secret := &coreV1.Secret{
        ObjectMeta: metaV1.ObjectMeta{
            Name:      secretname,
            Namespace: ns,
        },
        Type: coreV1.SecretTypeTLS,
        Data: data,
    }
	_, err := k.Client.CoreV1().Secrets(ns).Update(context.TODO(),secret,metaV1.UpdateOptions{})
    if err != nil {
        return err 
    }
	return nil 
}
func (k *K8sClient) CreateSecretByTlsCert(ns,secretname,cert,key string) error {
	defer k.CloseClient()
	data := map[string][]byte{
        coreV1.TLSCertKey:       []byte(cert),
        coreV1.TLSPrivateKeyKey: []byte(key),
    }
	secret := &coreV1.Secret{
        ObjectMeta: metaV1.ObjectMeta{
            Name:      secretname,
            Namespace: ns,
        },
        Type: coreV1.SecretTypeTLS,
        Data: data,
    }
	_, err := k.Client.CoreV1().Secrets(ns).Create(context.TODO(),secret,metaV1.CreateOptions{})
    if err != nil {
        return err 
    }
	return nil 
}
func (k *K8sClient) UpdateSecretByImageCert(ns,secretname,url,user,password string) error {
	defer k.CloseClient()

	dockerConfigJson := fmt.Sprintf(`{"auths":{"%v":{"username":"%v","password":"%v"}}}"`,url,user,password)
	secret := &coreV1.Secret{
        ObjectMeta: metaV1.ObjectMeta{
            Name: secretname,
        },
		Type: coreV1.SecretTypeDockerConfigJson,
        Data: map[string][]byte{
            ".dockerconfigjson": []byte(dockerConfigJson),
		},
    }
	_, err := k.Client.CoreV1().Secrets(ns).Update(context.TODO(), secret, metaV1.UpdateOptions{})
    if err != nil {
        return err 
    }
 
	return nil 
}
func (k *K8sClient) CreateSecretByImageCert(ns,secretname,url,user,password string) error {
	defer k.CloseClient()

	dockerConfigJson := fmt.Sprintf(`{"auths":{"%v":{"username":"%v","password":"%v"}}}"`,url,user,password)
	secret := &coreV1.Secret{
        ObjectMeta: metaV1.ObjectMeta{
            Name: secretname,
        },
		Type: coreV1.SecretTypeDockerConfigJson,
        Data: map[string][]byte{
            ".dockerconfigjson": []byte(dockerConfigJson),
		},
    }

    _, err := k.Client.CoreV1().Secrets(ns).Create(context.TODO(), secret, metaV1.CreateOptions{})
    if err != nil {
        return err 
    }
 
	return nil 
}
func (k *K8sClient) RoleBindingInfo(ns,name string) (*rbacV1.RoleBinding,error) {
	defer k.CloseClient()
	rolebinding,err := k.Client.RbacV1().RoleBindings(ns).Get(context.TODO(),name,metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return rolebinding,nil 
} 
func (k *K8sClient) NsInfo(ns string) (*coreV1.Namespace,error) {
	defer k.CloseClient()
	namespaceInfo,err := k.Client.CoreV1().Namespaces().Get(context.TODO(),ns,metaV1.GetOptions{})
	if err != nil {
		return namespaceInfo,err
	}
	return namespaceInfo,nil
}
func (k *K8sClient) RoleList(ns string) (*rbacV1.RoleList,error) {
	defer k.CloseClient()
	roles,err := k.Client.RbacV1().Roles(ns).List(context.TODO(),metaV1.ListOptions{})
	if err != nil {
		return nil,err 
	}
	return roles,nil 
}
func (k *K8sClient) RoleBindingList(ns string) (*rbacV1.RoleBindingList,error) {
	defer k.CloseClient()
	rolebindings,err := k.Client.RbacV1().RoleBindings(ns).List(context.TODO(),metaV1.ListOptions{})
	if err != nil {
		return nil,err 
	}
	return rolebindings,nil 
}
func (k *K8sClient) DeleteRoleBinding(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.RbacV1().RoleBindings(ns).Delete(context.TODO(),name,metaV1.DeleteOptions{})
	if err != nil {
		return err 
	}
	return nil 
}
func (k *K8sClient) NsList() (*coreV1.NamespaceList,error) {
	defer k.CloseClient()
	namespaceList,err := k.Client.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return namespaceList,err
	}
	return namespaceList,nil
}

func (k *K8sClient) PodList(ns string) (*coreV1.PodList,error) {
	var (
		podlist *coreV1.PodList = &coreV1.PodList{}
		err error 
	)
	defer k.CloseClient()
    // 列出Pods
	if ns != "" {
    	podlist, err = k.Client.CoreV1().Pods(ns).List(context.TODO(), metaV1.ListOptions{})
	} else {
		podlist, err = k.Client.CoreV1().Pods("").List(context.TODO(), metaV1.ListOptions{})
	}
	if err != nil {
        return nil,err 
    }
	return podlist,nil 
}
func (k *K8sClient) DeploymentList(ns string) (*appsV1.DeploymentList,error) {
	var (
		deployments *appsV1.DeploymentList = &appsV1.DeploymentList{}
		err error 
	)

	defer k.CloseClient()
	if ns != "" {
		deployments,err = k.Client.AppsV1().Deployments(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		deployments,err = k.Client.AppsV1().Deployments("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	deployments,nil 
}
func (k *K8sClient) SvcList(ns string) (*coreV1.ServiceList,error) {
	var (
		svcs *coreV1.ServiceList = &coreV1.ServiceList{}
		err error 
	)
	defer k.CloseClient()
	if ns != "" {
		svcs,err = k.Client.CoreV1().Services(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		svcs,err = k.Client.CoreV1().Services("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	svcs,nil 
}
func (k *K8sClient) StatefulSetList(ns string) (*appsV1.StatefulSetList,error) {
	var (
		statefulsets *appsV1.StatefulSetList = &appsV1.StatefulSetList{}
		err error 
	)
	defer k.CloseClient()
	if ns != "" {
		statefulsets,err = k.Client.AppsV1().StatefulSets(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		statefulsets,err = k.Client.AppsV1().StatefulSets("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	statefulsets,nil 
}
func (k *K8sClient) DaemonSetList(ns string) (*appsV1.DaemonSetList,error) {
	var (
		daemonsets *appsV1.DaemonSetList = &appsV1.DaemonSetList{}
		err error 
	)
	defer k.CloseClient()
	if ns != "" {
		daemonsets,err = k.Client.AppsV1().DaemonSets(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		daemonsets,err = k.Client.AppsV1().DaemonSets("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	daemonsets,nil 
} 
func (k *K8sClient) JobList(ns string) (*batchV1.JobList,error) {
	var (
		jobs *batchV1.JobList = &batchV1.JobList{}
		err error 
	)	
	defer k.CloseClient()
	if ns != "" {
		jobs,err = k.Client.BatchV1().Jobs(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		jobs,err = k.Client.BatchV1().Jobs("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	jobs,nil 
}
func (k *K8sClient) CronJobList(ns string) (*batchV1.CronJobList,error) {
	var (
		cronjobs *batchV1.CronJobList = &batchV1.CronJobList{}
		err error 
	)	
	defer k.CloseClient()
	if ns != "" {
		cronjobs,err = k.Client.BatchV1().CronJobs(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		cronjobs,err = k.Client.BatchV1().CronJobs("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	cronjobs,nil 
}
func (k *K8sClient) PodInfo(ns,podname string) (*coreV1.Pod,error) {
	defer k.CloseClient()
	pod, err := k.Client.CoreV1().Pods(ns).Get(context.TODO(), podname, metaV1.GetOptions{})
    if err != nil {
        return nil,err 
    }
	return pod,nil 
}

func (k *K8sClient) Event(ns,filter string) (*coreV1.EventList,error) {
	defer k.CloseClient()
	/*
	events, err := k.Client.CoreV1().Events(ns).List(context.TODO(), metaV1.ListOptions{
        FieldSelector: fmt.Sprintf("involvedObject.name=%s", podname),
    })
	*/
	events, err := k.Client.CoreV1().Events(ns).List(context.TODO(), metaV1.ListOptions{
        FieldSelector: filter,
    })
    if err != nil {
        return nil,err 
    }
	return events,nil 
}
func (k *K8sClient) NodeList() (*coreV1.NodeList,error) {
	defer k.CloseClient()
	nodes, err := k.Client.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
    if err != nil {
        return nil,err 
    }
	return nodes,nil 
}
func (k *K8sClient) IngressList(ns string) (*networkV1.IngressList,error) {
	var (
		ingresses *networkV1.IngressList = &networkV1.IngressList{}
		err error 
	)	
	defer k.CloseClient()
	if ns != "" {
		ingresses,err = k.Client.NetworkingV1().Ingresses(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		ingresses,err = k.Client.NetworkingV1().Ingresses("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	ingresses,nil 
}
func (k *K8sClient) ConfigMapList(ns string) (*coreV1.ConfigMapList,error) {
	var (
		configmaps *coreV1.ConfigMapList = &coreV1.ConfigMapList{}
		err error 
	)
	defer k.CloseClient()
	if ns != "" {
		configmaps,err = k.Client.CoreV1().ConfigMaps(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		configmaps,err = k.Client.CoreV1().ConfigMaps("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	configmaps,nil 
}
func (k *K8sClient) SecretList(ns string) (*coreV1.SecretList,error) {
	var (
		secrets *coreV1.SecretList = &coreV1.SecretList{}
		err error 
	)
	defer k.CloseClient()
	if ns != "" {
		secrets,err = k.Client.CoreV1().Secrets(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		secrets,err = k.Client.CoreV1().Secrets("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	secrets,nil 
} 
func (k *K8sClient) PvcList(ns string) (*coreV1.PersistentVolumeClaimList,error) {
	var (
		pvcs *coreV1.PersistentVolumeClaimList = &coreV1.PersistentVolumeClaimList{}
		err error 
	)
	defer k.CloseClient()
	if ns != "" {
		pvcs,err = k.Client.CoreV1().PersistentVolumeClaims(ns).List(context.Background(),metaV1.ListOptions{})
	} else {
		pvcs,err = k.Client.CoreV1().PersistentVolumeClaims("").List(context.Background(),metaV1.ListOptions{})
	}
	if err != nil {
		return nil,err 
	}
	return 	pvcs,nil 
}
func (k *K8sClient) PvList() (*coreV1.PersistentVolumeList,error) {
	var (
		pvs *coreV1.PersistentVolumeList = &coreV1.PersistentVolumeList{}
		err error 
	)
	defer k.CloseClient()
	pvs,err = k.Client.CoreV1().PersistentVolumes().List(context.Background(),metaV1.ListOptions{})
	if err != nil {
		return nil,err 
	}
	return 	pvs,nil 
}
func (k *K8sClient) StorageClassList() (*storageV1.StorageClassList,error) {
	var (
		storageclasses *storageV1.StorageClassList = &storageV1.StorageClassList{}
		err error 
	)
	storageclasses,err = k.Client.StorageV1().StorageClasses().List(context.Background(),metaV1.ListOptions{})
	if err != nil {
		return nil,err 
	}
	return 	storageclasses,nil 
}
func (k *K8sClient) NodeInfo(nodename string) (*coreV1.Node,error) {
	defer k.CloseClient()
	node,err := k.Client.CoreV1().Nodes().Get(context.TODO(),nodename,metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return node,nil 
}
func (k *K8sClient) RoleInfo(ns,role string) (*rbacV1.Role,error) {
	defer k.CloseClient()
	roleinfo,err := k.Client.RbacV1().Roles(ns).Get(context.TODO(),role,metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	} 
	return roleinfo,nil 
}
func (k *K8sClient) NodeLable(nodename string) (map[string]string,error) {
	defer k.CloseClient()
	node,err := k.Client.CoreV1().Nodes().Get(context.TODO(),nodename,metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return node.Labels,nil 
}

func (k *K8sClient) NodeTaint(nodename string) ([]coreV1.Taint,error) {
	defer k.CloseClient()
	node,err := k.Client.CoreV1().Nodes().Get(context.TODO(),nodename,metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return node.Spec.Taints,nil 
}

func (k *K8sClient) PatchNodeLable(nodename string,labels map[string]string) error {
	defer k.CloseClient()
	node,err := k.Client.CoreV1().Nodes().Get(context.TODO(),nodename,metaV1.GetOptions{})
	if err != nil {
		return err 
	}
	// 添加或更新标签
	if node.Labels == nil {
		node.Labels = map[string]string{}
	}
	for k,v := range labels {
		node.Labels[k] = v
	}
	// 更新节点
	_, err = k.Client.CoreV1().Nodes().Update(context.TODO(), node, metaV1.UpdateOptions{})
	if err != nil {
		return err 
	}
	return nil 
}
func (k *K8sClient) PatchNodeTaint(nodename string,taints map[string]string) error {
	defer k.CloseClient()
	node,err := k.Client.CoreV1().Nodes().Get(context.TODO(),nodename,metaV1.GetOptions{})
	if err != nil {
		return err 
	}
	
	for k,v := range taints {
		taint := &coreV1.Taint{
			Key:    k,
			Value:  v,
			Effect: coreV1.TaintEffectNoSchedule,
		}
		node.Spec.Taints = append(node.Spec.Taints, *taint)
	}
	// 更新节点
	_, err = k.Client.CoreV1().Nodes().Update(context.TODO(), node, metaV1.UpdateOptions{})
	if err != nil {
		return err 
	}
	return nil 
}

func (k *K8sClient) PatchNodeSchedule(nodename,schedulerule string) error {
	defer k.CloseClient()
	node,err := k.Client.CoreV1().Nodes().Get(context.TODO(),nodename,metaV1.GetOptions{})
	if err != nil {
		return err 
	}
	switch schedulerule {
	case "disable":
		node.Spec.Unschedulable = true
	case "enable":
		node.Spec.Unschedulable = false
	default:
		return errors.New("schedule rule invalid")
	}
	_, err = k.Client.CoreV1().Nodes().Update(context.TODO(), node, metaV1.UpdateOptions{})
    if err != nil {
        return err 
    }
	return nil 
}
func (k *K8sClient) DeploymentInfo(ns,deployname string) (*appsV1.Deployment,error) {
	defer k.CloseClient()
	deployinfo,err := k.Client.AppsV1().Deployments(ns).Get(context.Background(), deployname, metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return deployinfo,nil 
}
func (k *K8sClient) StatefulSetInfo(ns,statefulset string) (*appsV1.StatefulSet,error) {
	defer k.CloseClient()
	statefulsetinfo,err := k.Client.AppsV1().StatefulSets(ns).Get(context.Background(), statefulset, metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return statefulsetinfo,nil 
}
func (k *K8sClient) DaemonSetInfo(ns,daemonset string) (*appsV1.DaemonSet,error) {
	defer k.CloseClient()
	daemonsetinfo,err := k.Client.AppsV1().DaemonSets(ns).Get(context.Background(),daemonset,metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return daemonsetinfo,nil 
}
func (k *K8sClient) JobInfo(ns,job string) (*batchV1.Job,error) {
	defer k.CloseClient()
	jobinfo,err := k.Client.BatchV1().Jobs(ns).Get(context.Background(),job,metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return jobinfo,nil 
}
func (k *K8sClient) CronJobInfo(ns,cronjob string) (*batchV1.CronJob,error) {
	defer k.CloseClient()
	cronjobinfo,err := k.Client.BatchV1().CronJobs(ns).Get(context.Background(),cronjob,metaV1.GetOptions{})
	if err != nil {
		return nil,err 
	}
	return cronjobinfo,nil 
}
func (k *K8sClient) IngressInfo(ns,ingressname string) (*networkV1.Ingress,error) {
	defer k.CloseClient()
	ingress, err := k.Client.NetworkingV1().Ingresses(ns).Get(context.TODO(), ingressname, metaV1.GetOptions{})
    if err != nil {
        return nil,err 
    }
	return ingress,nil 
} 
func (k *K8sClient) SvcInfo(ns,svcname string)  (*coreV1.Service,error) {
	defer k.CloseClient()
	// 获取Service对象
	service, err := k.Client.CoreV1().Services(ns).Get(context.Background(), svcname, metaV1.GetOptions{})
	if err != nil {
		return nil,err  
	}
	return service,nil 
}
func (k *K8sClient) PatchNodeDrain(nodename string) error {
	defer k.CloseClient()
    pods, err := k.Client.CoreV1().Pods("").List(context.TODO(), metaV1.ListOptions{
        FieldSelector: "spec.nodeName=" + nodename,
    })
    if err != nil {
        return err 
    }
	for _, pod := range pods.Items {
        err := k.DeletePod(pod.Namespace,pod.Name)
		if err != nil {
        	return err 
        } 
    }
	return nil 
}
func (k *K8sClient) DeleteRole(ns,role string) error {
	defer k.CloseClient()
	err := k.Client.RbacV1().Roles(ns).Delete(context.TODO(),role,metaV1.DeleteOptions{})
	if err != nil {
		return err 
	}
	return nil 
}
func (k *K8sClient) DeleteServiceAccount(ns,serviceaccount string) error {
	defer k.CloseClient()
	deletePolicy := metaV1.DeletePropagationForeground
	err := k.Client.CoreV1().ServiceAccounts(ns).Delete(context.TODO(),serviceaccount,metaV1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err 
	}
	return nil 
}

func (k *K8sClient) DeletePod(ns,podname string) error {
	defer k.CloseClient()
	err := k.Client.CoreV1().Pods(ns).Delete(context.TODO(), podname, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil 
}
func (k *K8sClient) DeleteDeployment(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.AppsV1().Deployments(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil 
}
func (k *K8sClient) DeleteStatefulSet(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.AppsV1().StatefulSets(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil 
}
func (k *K8sClient) DeleteDaemonSet(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.AppsV1().DaemonSets(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil 
}
func (k *K8sClient) DeleteJob(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.BatchV1().Jobs(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil 
}
func (k *K8sClient) DeleteCronJob(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.BatchV1().CronJobs(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil 
}
func (k *K8sClient) DeleteService(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.CoreV1().Services(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil
}
func (k *K8sClient) DeleteIngress(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.NetworkingV1().Ingresses(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil
}
func (k *K8sClient) DeleteSecret(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.CoreV1().Secrets(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil
}
func (k *K8sClient) DeletePvc(ns,name string) error {
	defer k.CloseClient()
	err := k.Client.CoreV1().PersistentVolumeClaims(ns).Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil
}
func (k *K8sClient) DeletePv(name string) error {
	defer k.CloseClient()
	err := k.Client.CoreV1().PersistentVolumes().Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil
}
func (k *K8sClient) DeleteStorageClass(name string) error {
	defer k.CloseClient()
	err := k.Client.StorageV1().StorageClasses().Delete(context.TODO(), name, metaV1.DeleteOptions{GracePeriodSeconds: Int64ToPointInt64(0)})
	if err != nil {
		return err 
	} 
	return nil
}
func (k *K8sClient) PodsInNode(nodename string) (*coreV1.PodList,error) {
	defer k.CloseClient()
	pods, err := k.Client.CoreV1().Pods("").List(context.TODO(), metaV1.ListOptions{
        FieldSelector: "spec.nodeName=" + nodename,
    })
    if err != nil {
        return nil,err 
    }
	return pods,nil 
}
func (k *K8sClient) ConfigMapInfo(ns,mapname string) (*coreV1.ConfigMap,error) {
	defer k.CloseClient()
	mapinfo,err := k.Client.CoreV1().ConfigMaps(ns).Get(context.Background(),mapname,metaV1.GetOptions{})
	if err != nil {
        return nil,err 
    }
	return mapinfo,nil 
}
func (k *K8sClient) SecretInfo(ns,sectet string) (*coreV1.Secret,error) {
	defer k.CloseClient()
	secretinfo,err := k.Client.CoreV1().Secrets(ns).Get(context.Background(),sectet,metaV1.GetOptions{})
	if err != nil {
        return nil,err 
    }
	return secretinfo,nil 
}
func (k *K8sClient) PvcInfo(ns,pvcname string) (*coreV1.PersistentVolumeClaim,error) {
	defer k.CloseClient()
	pvcinfo,err := k.Client.CoreV1().PersistentVolumeClaims(ns).Get(context.Background(),pvcname,metaV1.GetOptions{})
	if err != nil {
        return nil,err 
    }
	return pvcinfo,nil
}
func (k *K8sClient) ServiceAccountInfo(ns,serviceaccount string) (*coreV1.ServiceAccount,error) {
	defer k.CloseClient()
	serviceaccountinfo,err := k.Client.CoreV1().ServiceAccounts(ns).Get(context.Background(),serviceaccount,metaV1.GetOptions{})
	if err != nil {
        return nil,err 
    }
	return serviceaccountinfo,nil
}
func (k *K8sClient) StorageClassInfo(ns,name string) (*storageV1.StorageClass,error) {
	defer k.CloseClient()
	storageClass, err := k.Client.StorageV1().StorageClasses().Get(context.TODO(), name, metaV1.GetOptions{})
    if err != nil {
        return nil,err 
    }
	return storageClass,nil 
}
func (k *K8sClient) PvInfo(pvname string) (*coreV1.PersistentVolume,error) {
	defer k.CloseClient()
	pvinfo,err := k.Client.CoreV1().PersistentVolumes().Get(context.Background(),pvname,metaV1.GetOptions{})
	if err != nil {
        return nil,err 
    }
	return pvinfo,nil
}
func (k *K8sClient) ClusterEvent() (*coreV1.EventList,error) {
	defer k.CloseClient()
	// 获取事件的watcher
    watcher, err := k.Client.CoreV1().Events("").List(context.Background(), metaV1.ListOptions{})
    if err != nil {
    	return nil,err 
    }
	return watcher,nil 
}
func (k *K8sClient) DeleteNode(node string) error {
	defer k.CloseClient()
	// 删除节点
    err := k.Client.CoreV1().Nodes().Delete(context.TODO(), node, metaV1.DeleteOptions{
        GracePeriodSeconds: Int64ToPointInt64(0), // 立即删除节点
    })
	if err != nil {
		return err 
	}
	return nil 
}
func  (k *K8sClient) DeployRollUpdate(deployment,ns string) error {
	defer k.CloseClient()
	patch := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"date": "%s"}}}}}`, time.Now().String())
    _, err := k.Client.AppsV1().Deployments(ns).Patch(context.TODO(), deployment, types.StrategicMergePatchType, []byte(patch), metaV1.PatchOptions{})
    if err != nil {
        return err 
    }
	return nil 
} 
func (k *K8sClient) DaemonSetRollUpdate(daemonset,ns string) error {
	defer k.CloseClient()
	dst,err := k.DaemonSetInfo(ns,daemonset)
	if err != nil {
		return err 
	}
	dst.Spec.UpdateStrategy = appsV1.DaemonSetUpdateStrategy{
        Type: appsV1.RollingUpdateDaemonSetStrategyType,
        RollingUpdate: &appsV1.RollingUpdateDaemonSet{
            MaxUnavailable: func(i intstr.IntOrString) *intstr.IntOrString { return &i }(intstr.FromInt(1)),
        },
    }
	_, err = k.Client.AppsV1().DaemonSets(ns).Update(context.TODO(), dst, metaV1.UpdateOptions{})
    if err != nil {
        return err 
    }
	return nil 
}
func (k *K8sClient) StatefulSetRollUpdate(statefulset,ns string) error {
	defer k.CloseClient()
	st,err := k.StatefulSetInfo(ns,statefulset)
	//st, err := k.Client.AppsV1().StatefulSets(ns).Get(context.TODO(), statefulset, metaV1.GetOptions{})
    if err != nil {
        return err 
    }
	st.Spec.Template.Labels["timestamp"] = fmt.Sprint(time.Now().Unix())
	_, err = k.Client.AppsV1().StatefulSets(ns).Update(context.TODO(), st, metaV1.UpdateOptions{})
    if err != nil {
        return err 
    }
	return nil 
}
func (k *K8sClient) Log(ns,podname string) (runtime.Object,error) {
	defer k.CloseClient()
	podLogs,err := k.Client.CoreV1().Pods(ns).GetLogs(podname, &coreV1.PodLogOptions{}).Do(context.TODO()).Get()
	if err != nil {
		return nil,err 
	}
	/*
    logsReq, err := http.NewRequest("GET", podLogs.URL(), nil)
    if err != nil {
        return "",err 
    }
	// 执行请求并读取日志
    logsResp, err := http.DefaultClient.Do(logsReq)
    if err != nil {
        return "",err 
    }
    defer logsResp.Body.Close()
 
    if logsResp.StatusCode != http.StatusOK {
        return "",err 
    }
 
    body, err := ioutil.ReadAll(logsResp.Body)
    if err != nil {
        return "",err 
    }
	*/
	return podLogs,nil 
}

func (k *K8sClient) CreateConfigMap(ns,configmap string,kv map[string]string) error {
	defer k.CloseClient()
	configMap := &coreV1.ConfigMap{
        ObjectMeta: metaV1.ObjectMeta{
            Name: configmap,
        },
        Data: kv,
    }
	_,err := k.Client.CoreV1().ConfigMaps(ns).Create(context.TODO(),configMap,metaV1.CreateOptions{})
    if err != nil {
        return err
    }
	return nil 
} 
func (k *K8sClient) ServiceAccountList() (*coreV1.ServiceAccountList,error) {
	defer k.CloseClient()
	serviceaccounts,err := k.Client.CoreV1().ServiceAccounts("").List(context.TODO(),metaV1.ListOptions{})
	if err != nil {
		return nil,err 
	}
	return serviceaccounts,nil 
}
func (k *K8sClient) UpdateConfigMap(ns,configmap string,kv map[string]string) error {
	defer k.CloseClient()
	cm,err := k.ConfigMapInfo(ns,configmap)
	if err != nil {
		return err 
	}
	cm.Data = kv
	_,err = k.Client.CoreV1().ConfigMaps(ns).Update(context.TODO(),cm,metaV1.UpdateOptions{})
    if err !=nil {
        return err
    }
    return nil
}
func (k *K8sClient) CreateResourceByYaml(obj interface{}) error {
	var err error 
	switch obj.(type) {
	case coreV1.Namespace:
		_,err = k.Client.CoreV1().Namespaces().Create(context.TODO(), obj.(*coreV1.Namespace),metaV1.CreateOptions{})
	case appsV1.Deployment:
		_,err = k.Client.AppsV1().Deployments("").Create(context.TODO(), obj.(*appsV1.Deployment),metaV1.CreateOptions{})
	case appsV1.StatefulSet:
		_,err = k.Client.AppsV1().StatefulSets("").Create(context.TODO(), obj.(*appsV1.StatefulSet),metaV1.CreateOptions{})
	case appsV1.DaemonSet:
		_,err = k.Client.AppsV1().DaemonSets("").Create(context.TODO(), obj.(*appsV1.DaemonSet),metaV1.CreateOptions{})
	case batchV1.Job:
		_,err = k.Client.BatchV1().Jobs("").Create(context.TODO(), obj.(*batchV1.Job),metaV1.CreateOptions{})
	case batchV1.CronJob:
		_,err = k.Client.BatchV1().CronJobs("").Create(context.TODO(), obj.(*batchV1.CronJob),metaV1.CreateOptions{})
	case coreV1.Pod:
		_,err = k.Client.CoreV1().Pods("").Create(context.TODO(), obj.(*coreV1.Pod),metaV1.CreateOptions{})
	case coreV1.Service:
		_,err = k.Client.CoreV1().Services("").Create(context.TODO(), obj.(*coreV1.Service),metaV1.CreateOptions{})
	case networkV1.Ingress:
		_,err = k.Client.NetworkingV1().Ingresses("").Create(context.TODO(), obj.(*networkV1.Ingress),metaV1.CreateOptions{})
	case coreV1.ConfigMap:
		_,err = k.Client.CoreV1().ConfigMaps("").Create(context.TODO(), obj.(*coreV1.ConfigMap),metaV1.CreateOptions{})
	case coreV1.Secret:
		_,err = k.Client.CoreV1().Secrets("").Create(context.TODO(), obj.(*coreV1.Secret),metaV1.CreateOptions{})
	case coreV1.PersistentVolumeClaim:
		_,err = k.Client.CoreV1().PersistentVolumeClaims("").Create(context.TODO(), obj.(*coreV1.PersistentVolumeClaim),metaV1.CreateOptions{})
	case coreV1.PersistentVolume:
		_,err = k.Client.CoreV1().PersistentVolumes().Create(context.TODO(), obj.(*coreV1.PersistentVolume),metaV1.CreateOptions{})
	case storageV1.StorageClass:
		_,err = k.Client.StorageV1().StorageClasses().Create(context.TODO(), obj.(*storageV1.StorageClass),metaV1.CreateOptions{})
	case coreV1.ServiceAccount:
		_,err = k.Client.CoreV1().ServiceAccounts("").Create(context.TODO(), obj.(*coreV1.ServiceAccount),metaV1.CreateOptions{})
	default:
		return errors.New("resource type invaild")
	}
	if err != nil {
		return err
	}

	return nil
}

func (k *K8sClient) CloseClient() {
	k.Client = nil
}
