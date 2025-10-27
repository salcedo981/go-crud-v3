package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_template_v3/pkg/config"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"
)

func AbsDiff(a, b int64) int64 {
	if a > b {
		return a - b
	}
	return b - a
}

func TimestampToUnix(timestamp string) int64 {
	t, _ := time.Parse(time.RFC3339, timestamp)
	return t.Unix()
}

func SendRequest(baseURL string, method string, body []byte, headers map[string]string, timeout int) (interface{}, error) {
	reqBody := bytes.NewBuffer(body)

	// Create the request
	req, err := http.NewRequest(method, baseURL, reqBody)
	if err != nil {
		return nil, err
	}

	// Set default content-type header if not provided
	if _, exists := headers["Content-Type"]; !exists {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	for key, value := range headers {
		fmt.Printf("HEADER: %s: %s\n", key, value)
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Handle empty response
	if len(body) == 0 {
		return nil, nil
	}

	// Try to parse response as JSON object
	var jsonRespObject map[string]interface{}
	if err := json.Unmarshal(body, &jsonRespObject); err == nil {
		return jsonRespObject, nil
	}

	// If parsing as JSON object fails, try as JSON array
	var jsonRespArray []interface{}
	if err := json.Unmarshal(body, &jsonRespArray); err == nil {
		return jsonRespArray, nil
	}

	// If neither parsing works, return an error
	return nil, fmt.Errorf("response is neither a JSON object nor a JSON array: %s", string(body))
}

func SpeakNotif(message string) {
	cmd := exec.Command("espeak", message)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

}

// Returns a pointer to val if non-zero, nil otherwise.
func Ptr[T comparable](val T) *T {
	var zero T
	if val == zero {
		return nil
	}
	return &val
}

func GetUserId(c fiber.Ctx) int {
	userIdInterface := c.Locals("userId")

	// Handle different possible types from JWT claims
	switch v := userIdInterface.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	default:
		return 0
	}
}

// Generic handler - just executes the query with the provided payload
func ExecuteDBFunction(c fiber.Ctx, query string, payload map[string]interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to marshal payload", err, http.StatusInternalServerError)
	}
	// Use a string variable to scan the result instead of []byte
	var resultStr string
	err = config.DBConnList[0].Debug().Raw(query, string(payloadJSON)).Scan(&resultStr).Error
	if err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "DB error", err, http.StatusInternalServerError)
	}

	// Convert string to []byte for JSON unmarshalling
	resultJSON := []byte(resultStr)
	var dbResponse map[string]interface{}
	if err := json.Unmarshal(resultJSON, &dbResponse); err != nil {
		return v1.JSONResponseWithError(c, respcode.ERR_CODE_500, "Failed to parse DB response", err, http.StatusInternalServerError)
	}

	code, _ := dbResponse["code"].(float64)
	codeInt := int(code)
	codeStr := strconv.Itoa(codeInt)
	message, _ := dbResponse["message"].(string)
	if message == ""{
		message = CodeMessageMap[codeStr]
	}

	if success, ok := dbResponse["success"].(bool); !ok || !success {
		return v1.JSONResponseWithError(c, codeStr, message, nil, codeInt)
	}

	data := dbResponse["data"]
	return v1.JSONResponseWithData(c, codeStr, message, data, codeInt)
}

// Returns map
func ExecuteDBFunctionRaw(dbFunc string, payload interface{}) (map[string]interface{}, error) {
    inputJSON, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal input: %w", err)
    }

    var resultStr string
    err = config.DBConnList[0].Raw((dbFunc), string(inputJSON)).Scan(&resultStr).Error
    if err != nil {
        return nil, fmt.Errorf("db error: %w", err)
    }

    var result map[string]interface{}
    if err := json.Unmarshal([]byte(resultStr), &result); err != nil {
        return nil, fmt.Errorf("failed to parse DB response: %w", err)
    }

    return result, nil
}


var CodeMessageMap = map[string]string{
	// Success codes
	respcode.SUC_CODE_200: respcode.SUC_CODE_200_MSG,
	respcode.SUC_CODE_201: respcode.SUC_CODE_201_MSG,
	respcode.SUC_CODE_202: respcode.SUC_CODE_202_MSG,
	respcode.SUC_CODE_203: respcode.SUC_CODE_203_MSG,
	respcode.SUC_CODE_204: respcode.SUC_CODE_204_MSG,
	respcode.SUC_CODE_205: respcode.SUC_CODE_205_MSG,
	respcode.SUC_CODE_206: respcode.SUC_CODE_206_MSG,

	// Error codes
	respcode.ERR_CODE_100:     respcode.ERR_CODE_100_MSG,
	respcode.ERR_CODE_100_CD:  respcode.ERR_CODE_100_CD_MSG,
	respcode.ERR_CODE_101_CD:  respcode.ERR_CODE_101_CD_MSG,
	respcode.ERR_CODE_102_CD:  respcode.ERR_CODE_102_CD_MSG,
	respcode.ERR_CODE_103_CD:  respcode.ERR_CODE_103_CD_MSG,
	respcode.ERR_CODE_104_CD:  respcode.ERR_CODE_104_CD_MSG,
	respcode.ERR_CODE_105_CD:  respcode.ERR_CODE_105_CD_MSG,
	respcode.ERR_CODE_106_CD:  respcode.ERR_CODE_106_CD_MSG,
	respcode.ERR_CODE_101:     respcode.ERR_CODE_101_MSG,
	respcode.ERR_CODE_102:     respcode.ERR_CODE_102_MSG,
	respcode.ERR_CODE_103:     respcode.ERR_CODE_103_MSG,
	respcode.ERR_CODE_104:     respcode.ERR_CODE_104_MSG,
	respcode.ERR_CODE_105:     respcode.ERR_CODE_105_MSG,
	respcode.ERR_CODE_106:     respcode.ERR_CODE_106_MSG,
	respcode.ERR_CODE_107:     respcode.ERR_CODE_107_MSG,
	respcode.ERR_CODE_108:     respcode.ERR_CODE_108_MSG,
	respcode.ERR_CODE_109:     respcode.ERR_CODE_109_MSG,
	respcode.ERR_CODE_110:     respcode.ERR_CODE_110_MSG,
	respcode.ERR_CODE_111:     respcode.ERR_CODE_111_MSG,
	respcode.ERR_CODE_111_MT:  respcode.ERR_CODE_111_MT_MSG,
	respcode.ERR_CODE_111_IT:  respcode.ERR_CODE_111_IT_MSG,
	respcode.ERR_CODE_112:     respcode.ERR_CODE_112_MSG,
	respcode.ERR_CODE_113:     respcode.ERR_CODE_113_MSG,
	respcode.ERR_CODE_114:     respcode.ERR_CODE_114_MSG,
	respcode.ERR_CODE_115:     respcode.ERR_CODE_115_MSG,
	respcode.ERR_CODE_116:     respcode.ERR_CODE_116_MSG,
	respcode.ERR_CODE_117:     respcode.ERR_CODE_117_MSG,
	respcode.ERR_CODE_118:     respcode.ERR_CODE_118_MSG,
	respcode.ERR_CODE_300:     respcode.ERR_CODE_300_MSG,
	respcode.ERR_CODE_301:     respcode.ERR_CODE_301_MSG,
	respcode.ERR_CODE_301_PR:  respcode.ERR_CODE_301_PR_MSG,
	respcode.ERR_CODE_302:     respcode.ERR_CODE_302_MSG,
	respcode.ERR_CODE_303:     respcode.ERR_CODE_303_MSG,
	respcode.ERR_CODE_304:     respcode.ERR_CODE_304_MSG,
	respcode.ERR_CODE_305:     respcode.ERR_CODE_305_MSG,
	respcode.ERR_CODE_306:     respcode.ERR_CODE_306_MSG,
	respcode.ERR_CODE_307:     respcode.ERR_CODE_307_MSG,
	respcode.ERR_CODE_308:     respcode.ERR_CODE_308_MSG,
	respcode.ERR_CODE_309:     respcode.ERR_CODE_309_MSG,
	respcode.ERR_CODE_310:     respcode.ERR_CODE_310_MSG,
	respcode.ERR_CODE_311:     respcode.ERR_CODE_311_MSG,
	respcode.ERR_CODE_312:     respcode.ERR_CODE_312_MSG,
	respcode.ERR_CODE_313:     respcode.ERR_CODE_313_MSG,
	respcode.ERR_CODE_314:     respcode.ERR_CODE_314_MSG,
	respcode.ERR_CODE_315:     respcode.ERR_CODE_315_MSG,
	respcode.ERR_CODE_316:     respcode.ERR_CODE_316_MSG,
	respcode.ERR_CODE_317:     respcode.ERR_CODE_317_MSG,
	respcode.ERR_CODE_318:     respcode.ERR_CODE_318_MSG,
	respcode.ERR_CODE_319:     respcode.ERR_CODE_319_MSG,
	respcode.ERR_CODE_330:     respcode.ERR_CODE_330_MSG,
	respcode.ERR_CODE_331:     respcode.ERR_CODE_331_MSG,
	respcode.ERR_CODE_400:     respcode.ERR_CODE_400_MSG,
	respcode.ERR_CODE_401:     respcode.ERR_CODE_401_MSG,
	respcode.ERR_CODE_402:     respcode.ERR_CODE_402_MSG,
	respcode.ERR_CODE_403:     respcode.ERR_CODE_403_MSG,
	respcode.ERR_CODE_404:     respcode.ERR_CODE_404_MSG,
	respcode.ERR_CODE_405:     respcode.ERR_CODE_405_MSG,
	respcode.ERR_CODE_406:     respcode.ERR_CODE_406_MSG,
	respcode.ERR_CODE_409:     respcode.ERR_CODE_409_MSG,
	respcode.ERR_CODE_500:     respcode.ERR_CODE_500_MSG,
	respcode.ERR_CODE_501:     respcode.ERR_CODE_501_MSG,
	respcode.ERR_CODE_502:     respcode.ERR_CODE_502_MSG,
}
