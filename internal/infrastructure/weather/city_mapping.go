package weather

import (
	"strings"
)

// CityMapping 城市名映射
type CityMapping struct {
	chineseToEnglish map[string]string
	englishToChinese map[string]string
}

// NewCityMapping 创建新的城市名映射
func NewCityMapping() *CityMapping {
	cm := &CityMapping{
		chineseToEnglish: make(map[string]string),
		englishToChinese: make(map[string]string),
	}

	// 初始化城市名映射
	cm.initCityMappings()

	return cm
}

// initCityMappings 初始化城市名映射
func (cm *CityMapping) initCityMappings() {
	// 主要城市映射
	mappings := map[string]string{
		// 直辖市
		"北京": "Beijing",
		"上海": "Shanghai",
		"天津": "Tianjin",
		"重庆": "Chongqing",

		// 省会城市
		"广州":   "Guangzhou",
		"深圳":   "Shenzhen",
		"杭州":   "Hangzhou",
		"南京":   "Nanjing",
		"成都":   "Chengdu",
		"武汉":   "Wuhan",
		"西安":   "Xian",
		"济南":   "Jinan",
		"青岛":   "Qingdao",
		"大连":   "Dalian",
		"沈阳":   "Shenyang",
		"哈尔滨":  "Harbin",
		"长春":   "Changchun",
		"石家庄":  "Shijiazhuang",
		"太原":   "Taiyuan",
		"呼和浩特": "Hohhot",
		"郑州":   "Zhengzhou",
		"合肥":   "Hefei",
		"南昌":   "Nanchang",
		"福州":   "Fuzhou",
		"厦门":   "Xiamen",
		"长沙":   "Changsha",
		"南宁":   "Nanning",
		"海口":   "Haikou",
		"贵阳":   "Guiyang",
		"昆明":   "Kunming",
		"拉萨":   "Lhasa",
		"兰州":   "Lanzhou",
		"西宁":   "Xining",
		"银川":   "Yinchuan",
		"乌鲁木齐": "Urumqi",

		// 其他重要城市
		"苏州": "Suzhou",
		"无锡": "Wuxi",
		"宁波": "Ningbo",
		"温州": "Wenzhou",
		"佛山": "Foshan",
		"东莞": "Dongguan",
		"中山": "Zhongshan",
		"珠海": "Zhuhai",
		"惠州": "Huizhou",
		"江门": "Jiangmen",
		"肇庆": "Zhaoqing",
		"清远": "Qingyuan",
		"韶关": "Shaoguan",
		"河源": "Heyuan",
		"梅州": "Meizhou",
		"汕尾": "Shanwei",
		"阳江": "Yangjiang",
		"茂名": "Maoming",
		"湛江": "Zhanjiang",
		"潮州": "Chaozhou",
		"揭阳": "Jieyang",
		"云浮": "Yunfu",

		// 区级地名映射（部分）
		"北京海淀":  "Beijing",
		"北京朝阳":  "Beijing",
		"北京西城":  "Beijing",
		"北京东城":  "Beijing",
		"北京丰台":  "Beijing",
		"北京石景山": "Beijing",
		"北京门头沟": "Beijing",
		"北京房山":  "Beijing",
		"北京通州":  "Beijing",
		"北京顺义":  "Beijing",
		"北京昌平":  "Beijing",
		"北京大兴":  "Beijing",
		"北京怀柔":  "Beijing",
		"北京平谷":  "Beijing",
		"北京密云":  "Beijing",
		"北京延庆":  "Beijing",

		"上海浦东": "Shanghai",
		"上海黄浦": "Shanghai",
		"上海徐汇": "Shanghai",
		"上海长宁": "Shanghai",
		"上海静安": "Shanghai",
		"上海普陀": "Shanghai",
		"上海虹口": "Shanghai",
		"上海杨浦": "Shanghai",
		"上海闵行": "Shanghai",
		"上海宝山": "Shanghai",
		"上海嘉定": "Shanghai",
		"上海金山": "Shanghai",
		"上海松江": "Shanghai",
		"上海青浦": "Shanghai",
		"上海奉贤": "Shanghai",
		"上海崇明": "Shanghai",
	}

	// 填充映射表
	for chinese, english := range mappings {
		cm.chineseToEnglish[chinese] = english
		cm.englishToChinese[english] = chinese
	}
}

// GetEnglishName 获取英文城市名
func (cm *CityMapping) GetEnglishName(chineseName string) (string, bool) {
	// 直接匹配
	if english, exists := cm.chineseToEnglish[chineseName]; exists {
		return english, true
	}

	// 模糊匹配：检查是否包含已知的中文城市名
	for chinese, english := range cm.chineseToEnglish {
		if strings.Contains(chineseName, chinese) {
			return english, true
		}
	}

	return "", false
}

// GetChineseName 获取中文城市名
func (cm *CityMapping) GetChineseName(englishName string) (string, bool) {
	chinese, exists := cm.englishToChinese[englishName]
	return chinese, exists
}

// IsChineseCity 判断是否为中文城市名
func (cm *CityMapping) IsChineseCity(cityName string) bool {
	// 简单的判断：包含中文字符
	for _, char := range cityName {
		if char >= 0x4e00 && char <= 0x9fff {
			return true
		}
	}
	return false
}
