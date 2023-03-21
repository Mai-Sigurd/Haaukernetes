package apis

type Challenge struct {
	Port      int64  `json:"port"`
	ImageName string `json:"imageName"`
	Namespace string `json:"namespace"`
}
