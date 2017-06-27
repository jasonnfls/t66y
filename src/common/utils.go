package common

import (
    "io/ioutil"
    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/transform"
    "bytes"
)

func Decodegbk(s []byte) ([]byte, error) {
    I := bytes.NewReader(s)
    O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
    d, e := ioutil.ReadAll(O)
    if e != nil {
        return nil, e
    }
    return d, nil
}
func DecodegbkStr(s string) string {
    I := bytes.NewReader([]byte(s))
    O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
    d, e := ioutil.ReadAll(O)
    if e != nil {
        panic(e)
    }
    return string(d)
}

