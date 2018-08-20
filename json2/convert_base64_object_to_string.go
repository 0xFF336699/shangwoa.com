package json2

// ConvertBase64ObjectToString 在数据库返回数据不使用struct而用map的时候会把数据库里的数组或对象被改为base64转码后的对象，这样在把这个对象再次解析为struct的时候会因为输入的是对象，而不是字符串（例如pq.StringArray要求输入的是字符串才能解析，好像byte也行），而实际输入的是json的对象结构，就导致无法解析了。
//s := "select " + strings.Join(cols, ",") + " FROM " + "shortcode_media_vo" + " WHERE id = $1"
//ok, row, err := gosqlmf.QueryOne(db, s, id)
// row 返回的就是{"id":4496,"tags":"e30="} e30=解析出来的是{}，而不是pq.Stringarray需要的"{}"，两边缺少引号所以被认为是对象。
// ConvertBase64ObjectToString(row) 将会把 e30=转为字符串的{}
func ConvertBase64ObjectToString(m map[string]interface{}) {
	for name, value := range m {
		if u8, ok := value.([]uint8); ok {
			m[name] = string(u8)
		}
	}
}
