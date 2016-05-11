package file

import (
	. "github.com/fishedee/util"
	. "github.com/fishedee/web"
)

type uploadAoTest struct {
	Test
	UploadAo UploadAoModel
}

func (this *uploadAoTest) TestFile() {
	testCase := []string{
		"",
		"helloworld",
		"你好 dd",
	}
	for _, singleTestCase := range testCase {
		url := this.UploadAo.UploadFile([]byte(singleTestCase))
		var data string
		err := DefaultAjaxPool.Get(&Ajax{
			Url:          url,
			ResponseData: &data,
		})
		this.AssertEqual(err == nil, true, singleTestCase)
		this.AssertEqual(data, singleTestCase, singleTestCase)
	}
}
