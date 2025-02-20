package engine

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"

	"github.com/dop251/goja"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"

	"kd_api/internal/spider/goreq"
)

func encryptLogonParams(username, password, sessionId, desKey, code string) (string, error) {
	vm := goja.New()
	jsScript, err := os.ReadFile(jsFilename)
	if err != nil {
		return "", err
	}
	_, err = vm.RunString(string(jsScript))
	if err != nil {
		return "", err
	}
	jsEncrypt, ok := goja.AssertFunction(vm.Get("getParams"))
	if !ok {
		return "", errors.New("goja: Run js error")
	}
	key, err := jsEncrypt(
		goja.Undefined(), vm.ToValue(username), vm.ToValue(password),
		vm.ToValue(sessionId), vm.ToValue(desKey), vm.ToValue(code))
	if err != nil {
		return "", err
	}
	return goreq.QueriesEscape(key.String()), err
}

func encryptClassSchedule(xn, xq, xh string) string {
	params := "xn=" + xn + "&xq=" + xq + "&xh=" + xh
	paramsBytes := []byte(params)
	return base64.StdEncoding.EncodeToString(paramsBytes)
}

func (c *KdClient) getValidateCode() (string, error) {
	p := "37"
	api := fmt.Sprintf("%s?v=%s", validateCodeApi, p)
	res, err := c.Get(api, defaultHeaders, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New(res.Status)
	}
	// OpenCV
	imgData, err := io.ReadAll(res.Body)
	img, err := gocv.IMDecode(imgData, gocv.IMReadColor)
	if img.Empty() {
		return "", errors.New("gocv.IMRead Error")
	}
	defer img.Close()
	// Cov Gray
	grey := gocv.NewMat()
	defer grey.Close()
	gocv.CvtColor(img, &grey, gocv.ColorBGRToGray)
	// Binary
	binary := gocv.NewMat()
	defer binary.Close()
	gocv.Threshold(grey, &binary, 128, 255, gocv.ThresholdBinary)
	resized := gocv.NewMat()
	defer resized.Close()
	point := image.Point{X: binary.Cols() * 3, Y: binary.Rows() * 3}
	gocv.Resize(binary, &resized, point, 0, 0, gocv.InterpolationLinear)
	codePath := tmpDir + "/code.jpeg"
	gocv.IMWrite(codePath, resized)
	newCodePath := tmpDir + "/code_.jpeg"
	err = exec.Command("convert", codePath, "-units", "PixelsPerInch", "-density", "300", newCodePath).Run()
	// OCR
	cli := gosseract.NewClient()
	defer cli.Close()
	err = cli.SetImage(newCodePath)
	if err != nil {
		return "", err
	}
	err = cli.SetLanguage("eng")
	err = cli.SetPageSegMode(gosseract.PSM_SINGLE_LINE)
	err = cli.SetVariable("tessedit_char_whitelist", "0123456789")
	if err != nil {
		return "", err
	}
	text, err := cli.Text()
	if err != nil {
		return "", err
	}
	return text, nil
}
