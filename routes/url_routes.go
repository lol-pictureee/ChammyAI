package routes

type route string

// FS - FACE SWAP
const (
	ApiFsNiceFish route = "https://aifaceswapper.io/api/nicefish/fs/singleface" // main api FaceSwapper service
	ApiFsSegmind  route = "https://api.segmind.com/v1/faceswap-v3"              // second reserve api FaceSwapper service
	ApiFuseBrain  route = "https://api-key.fusionbrain.ai/"                     // Кандинский АПИ
	ApiUrlUpload  route = "https://api.imageban.ru/v1/upload"                   // Апи загрузки URL для последующей передачи IMAGE BAN API
	ApiQuotes     route = "https://citaty.info/random"
)

func (r route) String() string {
	return string(r)
}
