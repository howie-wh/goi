package config

import (
	"encoding/json"
	"io/ioutil"

	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/model"

	log "github.com/sirupsen/logrus"
)

var defaultSpexInject = `
{
    "import_configs": [
        {
            "type": "import",
            "file": "internal/spexclient/spexclient.go",
            "code": "import \"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/mock\""
        }
    ],
    "function_configs": [
        {
            "type": "function",
            "file": "internal/spexclient/spexclient.go",
            "relative_line": 0,
            "func": "RPCRequest",
            "code": "if code, isIgnoreMock := mock.MockerForward(ctx, command, request, response); !isIgnoreMock { return code };defer func() { _ = mock.RecorderForward(ctx, command, request, response) }()"
        }
    ]
}
`
var defaultHTTPInject = `
{
    "import_configs": [
        {
            "type": "import",
            "file": "internal/decorator/http/http.go",
            "code": "import \"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/mock\""
        }
    ],
    "function_configs": [
        {
            "type": "function",
            "file": "internal/decorator/http/http.go",
            "relative_line": 0,
            "func": "GetJSON",
            "code": "if code, isIgnoreMock := mock.MockerForward(ctx, url, params, respData); !isIgnoreMock { return mock.ErrorTranslate(code)};defer func() { _ = mock.RecorderForward(ctx, url, params, respData) }()"
        },
        {
            "type": "function",
            "file": "internal/decorator/http/http.go",
            "relative_line": 0,
            "func": "PostJSON",
            "code": "if code, isIgnoreMock := mock.MockerForward(ctx, url, data, respData); !isIgnoreMock { return mock.ErrorTranslate(code) };defer func() { _ = mock.RecorderForward(ctx, url, data, respData) }()"
        }
    ]
}
`
var defaultCacheInject = `
{
    "import_configs": [
        {
            "type": "import",
            "file": "internal/platform/cache/unified_cache_client.go",
            "code": "import \"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/mock\""
        }
    ],
    "function_configs": [
        {
            "type": "function",
            "file": "internal/platform/cache/unified_cache_client.go",
            "relative_line": 0,
            "func": "GetWithContext",
            "code": "if code, isIgnoreMock := mock.MockerForward(ctx, key, \"\", value); !isIgnoreMock { return mock.ErrorTranslate(code)}"
        }
    ]
}
`

// GetDefaultSpexInjectConfig ...
func GetDefaultSpexInjectConfig() *model.InjectConfig {
	spexInjectConfig := &model.InjectConfig{}
	err := json.Unmarshal([]byte(defaultSpexInject), spexInjectConfig)
	if err != nil {
		log.Fatalf("Umarshal failed: %v", err)
		return nil
	}
	log.Debug(spexInjectConfig)
	return spexInjectConfig
}

// GetDefaultHTTPInjectConfig ...
func GetDefaultHTTPInjectConfig() *model.InjectConfig {
	httpInjectConfig := &model.InjectConfig{}
	err := json.Unmarshal([]byte(defaultHTTPInject), httpInjectConfig)
	if err != nil {
		log.Fatalf("Umarshal failed: %v", err)
	}
	return httpInjectConfig
}

// GetDefaultCacheInjectConfig ...
func GetDefaultCacheInjectConfig() *model.InjectConfig {
	cacheInjectConfig := &model.InjectConfig{}
	err := json.Unmarshal([]byte(defaultCacheInject), cacheInjectConfig)
	if err != nil {
		log.Fatalf("Umarshal failed: %v", err)
	}
	return cacheInjectConfig
}

// GetCustomInjectConfig ...
func GetCustomInjectConfig(filename string) *model.InjectConfig {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Read file %s: %s", filename, err)
	}
	log.Debug(string(fileBytes))

	customInjectConfig := &model.InjectConfig{}
	err = json.Unmarshal(fileBytes, customInjectConfig)
	if err != nil {
		log.Fatalf("Umarshal failed: %v", err)
	}
	log.Debug(customInjectConfig)
	return customInjectConfig
}

// MergeInjectConfig ...
func MergeInjectConfig(target *model.InjectConfig, source *model.InjectConfig) {
	target.ImportConfigs = append(target.ImportConfigs, source.ImportConfigs...)
	target.FunctionConfigs = append(target.FunctionConfigs, source.FunctionConfigs...)
}
