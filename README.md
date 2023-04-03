DataTester GO SDK
==================

## Limitation
>This SDK is only supported on **Go v1.14** and later versions.

## Prerequisite
Obtain your project's App Key:
1. Go to the [BytePlus console](https://console.byteplus.com/) and sign in to your account.
2. Under the **Products** section, click **BytePlus Data Intelligence**.
3. On the **Project List** page, for the project you want to integrate the SDK with, under the **Actions** column, click **Details**.
4. On the **Social Media Details** pop-up window, copy the **App Key**.

## Install the SDK
```
go get github.com/volcengine/datatester-go-sdk@v1.0.4
```

## Using the SDK
The following is a code example of using the Go SDK.
```
package main

import (
	"github.com/volcengine/datatester-go-sdk/client"
	"github.com/volcengine/datatester-go-sdk/config"
)

func main() {
	// Find your app key by clicking "details" on the project list page on the BytePlus console
	abClient := client.NewClient("${app_key}")

    /*
	   client.NewClient("2b47a1f318d78fd718548153901addde",
	   config.WithMetaHost(config.MetaHostSG), // Change it to "https://datarangers.com"
	   config.WithTrackHost(config.TrackHostSG), // Change it to "https://mcs.tobsnssdk.com"
	   config.WithWorkerNumOnce(20)) // The number of worker threads. 20 by default. Only supports one-time configuration
	   config.WithFetchInterval(60 * time.Second), // meta update interval, N * time.Second is recommended
       config.WithAnonymousConfig(true, true), // anonymous user event reporting configuration(enable/disable, saas/onpremise)
       config.WithLogger(log.NewLogrusAdapt(logrus.New()))) // customize the log interface
	*/

	// attributes: User properties. Only used for allocating the traffic
	attributes := map[string]interface{}{
	}
	// decisionId: The local user identifier. Not used for event tracking. You need to replace it with the actual user ID
	// trackId(uuid): The user identifier for event tracking. You need to replace it with the actual user ID
	value, err := abClient.Activate("${experiment_key}", "decisionId", "trackId", true, attributes)
	if err != nil {
		return
	}
	if value.(bool) {

	} else {

	}
}
```

## API reference
>Interfaces with **WithImpression** automatically track events. Meanwhile, make sure you fill in the **trackId** field in the **activate** interface if you want to track events.

### AbClient
A class for traffic allocation during initialization.
```
NewClient(token string, configs ...config.Func) *AbClient
```
#### Parameter description

| Parameter                                            | Description                                                                       | Value                            |
|:-----------------------------------------------------|:----------------------------------------------------------------------------------|:---------------------------------|
| token                                                | Your project's App Key. You can obtain it in the Prerequisites section.           | 2b47*****8d78fd718548153901addde |
| config.WithMetaHost(config.MetaHostSG)               | Default value is CN SAAS, set according to business needs                         |                                  |
| config.WithFetchInterval(60 * time.Second)           | Meta update interval, N * time.Second is recommended or use default value         |                                  |
| config.WithTrackHost(config.TrackHostSG)             | Default value is CN SAAS, set according to business needs                         |                                  |
| config.WithWorkerNumOnce(20))                        | The number of worker threads. 20 by default. Only supports one-time configuration |                                  |
| config.WithAnonymousConfig(true, true)               | Anonymous user event reporting configuration(enable/disable, saas/onpremise)      |                                  |
| config.WithLogger(log.NewLogrusAdapt(logrus.New()))) | The logger interface which has a default value but you can customize it.          |                                  |

### NewClientWithUserAbInfo
A class for traffic allocation during initialization.
```
NewClientWithUserAbInfo(token string, userAbInfoHandler handler.UserAbInfoHandler, configs ...config.Func) *AbClient
```
#### Parameter description

| Parameter         | Description                                                                                                                                                                      | Value |
|:------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:------|
| userAbInfoHandler | Ensure that incoming users' information about version ID remains unchanged. If you want to constantly store the information of incoming users, you must implement it by yourself |       |

### Activate
Obtain a specific experiment version's configuration after traffic allocation and automatically track the exposed events.
```
func (t *AbClient) Activate(variantKey, decisionId, trackId string, defaultValue interface{},
	attributes map[string]interface{}) (interface{}, error)
```
#### Parameter description

| Parameter     | Description                                                                                                             |
|:--------------|:------------------------------------------------------------------------------------------------------------------------|
| variantKey    | The key of the experiment version                                                                                       |
| decisionId    | The local user identifier for traffic allocating                                                                        |
| trackId       | The user identifier for event tracking. You need to replace it with the actual user ID                                  |
| defaultValue	 | When a user/device is not in this version, then the value of this parameter is returned. You can set its value as "nil" |
| attributes    | User properties                                                                                                         |

#### Returned value
Parameters of the version ID that a user enters or the default value when a user/device is not in this version.

### getExperimentVariantName
Obtain the version's name of an experiment that a user enters.
```
func (t *AbClient) GetExperimentVariantName(experimentId, decisionId string,
	attributes map[string]interface{}) (string, error)
```
#### Parameter description
| Parameter    | Description                                     |
|:-------------|:------------------------------------------------|
| experimentId | The experiment ID to which traffic is allocated |

#### Returned value
Name of the version ID that a user enters or "" when a user/device is not in this version.

### getExperimentConfigs
Obtain detailed information about the version that a user enters.
```
func (t *AbClient) GetExperimentConfigs(experimentId, decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### Parameter description

#### Returned value
The detailed information about the version ID that a user enters or "nil" when a user/device is not in this version. The following is an example:
```
{
    "father_code": {
       "val": "father_code_2",
       "vid": "12345"
    }
}
```

### getAllExperimentConfigs
Obtain detailed information about all the versions in all experiments.
```
func (t *AbClient) GetAllExperimentConfigs(decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### Parameter description

#### Returned value
The detailed information about version IDs that a user enters or "nil" when a user/device is not in this version. The following is an example:
```
{
    "father_code": {
       "val": "father_code_2",
       "vid": "12345"
    },
    "grey_rollout": {
        "val": false,
        "vid": "45678"
    }
}
```

### getFeatureConfigs
Obtain detailed information about a feature that a user joins.
```
func (t *AbClient) GetFeatureConfigs(featureId, decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### Parameter description
| Parameter | Description              |
|:----------|:-------------------------|
| featureId | The feature's identifier |

#### Returned value
An array object of detailed information about a feature's variant that a user joins. "nil" when a user/device is not in this feature or the feature is disabled. The following is an example:
```
{
   "feature_key":{
        "val" : "prod",
        "vid" : "20006421"
    }
}
```

### getAllFeatureConfigs
Obtain detailed information about all features that a user joins.
```
func (t *AbClient) GetAllFeatureConfigs(decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### Parameter description

#### Returned value
The detailed information about all feature variants that a user joins. The following is an example:
```
{
   "feature_key":{
       "val" : "prod",
       "vid" : "20006421"
    }
    "feature_key_color":{
       "val" : "true",
       "vid" : "20006423"
    }
}
```

### getExperimentVariantNameWithImpression
Obtain the version's name of the experiment that a user enters.
```
func (t *AbClient) GetExperimentVariantNameWithImpression(experimentId, decisionId, trackId string,
	attributes map[string]interface{}) (string, error)
```
#### Parameter description

#### Returned value
Name of the version ID that a user enters or "" when a user/device is not in this version.

### getExperimentConfigsWithImpression
Obtain detailed information about the version that a user enters.
```
func (t *AbClient) GetExperimentConfigsWithImpression(experimentId, decisionId, trackId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### Parameter description

#### Returned value
The detailed information about the version ID that a user enters or "nil" when a user/device is not in this version. The following is an example:
```
{
    "father_code": {
       "val": "father_code_2",
       "vid": "12345"
    }
}
```

### getAllExperimentConfigsWithImpression
Obtain detailed information about all the versions in all experiments.
```
func (t *AbClient) GetAllExperimentConfigsWithImpression(decisionId string, trackId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### Parameter description

#### Returned value
The detailed information about version IDs that a user enters or "nil" when a user/device is not in this version. The following is an example:
```
{
    "father_code": {
       "val": "father_code_2",
       "vid": "12345"
    },
    "grey_rollout": {
        "val": false,
        "vid": "45678"
    }
}
```

### getFeatureConfigsWithImpression
Obtain detailed information about a feature that a user joins.
```
func (t *AbClient) GetFeatureConfigsWithImpression(featureId, decisionId, trackId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### Parameter description

#### Returned value
The detailed information about a feature's variant that a user joins. "nil" when a user/device is not in this feature or the feature is disabled. The following is an example:
```
{
   "feature_key":{
        "val" : "prod",
        "vid" : "20006421"
    }
}
```

## Others
In order to better use the sdk, some suggestions are provided.

### UserAbInfoHandler
>Maintain user historical decision information; If you need to use the function of '**freezing experiment**' or '**Traffic changes will not affect exposed users**', you can customize the implementation class processing, and use it in when instantiating AbClient(NewClientWithUserAbInfo).
#### Suggestion
>It is recommended to use redis to cache decision information.

The following is an example:
```
client.NewClientWithUserAbInfo("appKey", NewRedisUserAbInfoHandler())

type RedisAbInfoHandler struct{}

func (u *RedisAbInfoHandler) Query(decisionID string) string {
	// need to implement it yourself
    return redis.get(decisionID);
}

func (u *RedisAbInfoHandler) CreateOrUpdate(decisionID, experiment2Variant string) bool {
    // need to implement it yourself
	return redis.set(decisionID, experiment2Variant);
}

func (u *RedisAbInfoHandler) NeedPersistData() bool {
    // return true if customize this interface
	return true
}

func NewRedisUserAbInfoHandler() *RedisAbInfoHandler {
	return &RedisAbInfoHandler{}
}
```

### Anonymously tracking
>If there is no uuid as trackId, you can use device_id, web_id, bddid(only onpremise) for anonymous tracking.
1. init anonymous config
```
enable anonymously tracking
client.NewClient("appKey", config.WithAnonymousConfig(true, true))
```
2. append device_id, web_id, bddid to attributes when trackId is empty string
```
trackId = "";
attributes["device_id"] = 1234; int64
attributes["web_id"] = 5678; int64
attributes["bddid"] = "91011"; string
```
3. call activate or 'WithImpression' interfaces


DataTester GO SDK
==================

## 版本需求
>**Go v1.14**及更高版本

## 准备工作
获取应用的App Key（即SDK使用的token）:
1. 访问[火山引擎](https://www.volcengine.com/)并登录您的账号
2. 进入集团设置页面，找到应用列表-应用ID列
3. 鼠标悬停在应用ID后的感叹号上获取App Key

## 安装SDK
```
go get github.com/volcengine/datatester-go-sdk@v1.0.4
```

## 代码示例
```
package main

import (
	"github.com/volcengine/datatester-go-sdk/client"
	"github.com/volcengine/datatester-go-sdk/config"
)

func main() {
	abClient := client.NewClient("${app_key}")

    /*
	   client.NewClient("2b47a1f318d78fd718548153901addde",
	   config.WithMetaHost(config.MetaHostCN), // 默认使用国内SAAS域名，私有化需要自行传入产品域名
	   config.WithTrackHost(config.TrackHostCN), // 默认使用国内SAAS域名，私有化需要自行传入上报域名
	   config.WithWorkerNumOnce(20)) // 事件上报协程数，一般不需要设置
	   config.WithFetchInterval(60 * time.Second), // meta更新间隔，默认为60s，一般不需要设置
       config.WithAnonymousConfig(true, true), // 匿名上报配置，第一个参数为开启关闭，第二个参数区分saas和私有化
       config.WithLogger(log.NewLogrusAdapt(logrus.New()))) // 自定义日志接口，提供默认实现
	*/

	// attributes: 用户属性
	attributes := map[string]interface{}{
	}
	// decisionId: 本地分流用户标识，不用于事件上报，请替换为客户的真实用户标识
	// trackId(uuid): 事件上报用户标识，用于事件上报，请替换为客户的真实用户标识
	value, err := abClient.Activate("${experiment_key}", "decisionId", "trackId", true, attributes)
	if err != nil {
		return
	}
	if value.(bool) {

	} else {

	}
}
```

## 接口描述

### AbClient
初始化ABTest分流类
```
NewClient(token string, configs ...config.Func) *AbClient
```
#### 参数

| 参数                                                    | 描述                                         | 值                          |
|:------------------------------------------------------|:-------------------------------------------|:---------------------------|
| token                                                 | 获取到的App Key                                | 2b47*****8d78fd71854815390 |
| config.WithMetaHost                                   | 默认使用国内SAAS的域名，根据业务需要设置                     |                            |
| config.WithFetchInterval(60 * time.Second)            | meta更新间隔，默认为60s，一般不需要设置                    |                            |
| config.WithTrackHost                                  | 默认使用国内SAAS的域名，根据业务需要设置                     |                            |
| config.WithWorkerNumOnce(20))                         | 事件上报协程数，一般不需要设置                            |                            |
| config.WithAnonymousConfig(true, true)                | 匿名上报配置，第一个参数为开启关闭，第二个参数区分saas和私有化          |                            |
| config.WithLogger(log.NewLogrusAdapt(logrus.New())))  | 自定义日志接口，提供默认实现                             |                            |

### NewClientWithUserAbInfo
初始化ABTest分流类，传入自定义的userAbInfoHandler
```
NewClientWithUserAbInfo(token string, userAbInfoHandler handler.UserAbInfoHandler, configs ...config.Func) *AbClient
```
#### 参数

| 参数                | 描述                                   | 值   |
|:------------------|:-------------------------------------|:----|
| userAbInfoHandler | 用户进组信息管理接口，提供默认实现，实验冻结和进组不出组场景下需自行实现 |     |

### Activate
获取特定key的分流结果，并上报曝光事件
```
func (t *AbClient) Activate(variantKey, decisionId, trackId string, defaultValue interface{},
	attributes map[string]interface{}) (interface{}, error)
```
#### 参数

| 参数            | 描述       |
|:--------------|:---------|
| variantKey    | 变体的key   |
| decisionId    | 本地分流用户标识 |
| trackId       | 事件上报用户标识 |
| defaultValue	 | 变体默认值    |
| attributes    | 用户属性     |

#### 返回值
该函数返回命中版本的参数值，未命中时返回默认值

### getExperimentVariantName
获取用户命中的特定实验的变体名称
```
func (t *AbClient) GetExperimentVariantName(experimentId, decisionId string,
	attributes map[string]interface{}) (string, error)
```
#### 参数
| 参数           | 描述        |
|:-------------|:----------|
| experimentId | 指定分流的实验Id |

#### 返回值
该函数返回用户命中的特定实验的变体名称

### getExperimentConfigs
获取用户命中的特定实验的变体详情
```
func (t *AbClient) GetExperimentConfigs(experimentId, decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### 参数

#### 返回值
该函数返回用户命中某个实验的变体详情，通常仅能命中一个变体
```
{
    "father_code": {
       "val": "father_code_2",
       "vid": "12345"
    }
}
```

### getAllExperimentConfigs
获取用户命中的所有实验的变体详情
```
func (t *AbClient) GetAllExperimentConfigs(decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### 参数

#### 返回值
该函数返回用户命中所有实验的变体详情，通常命中多个变体
```
{
    "father_code": {
       "val": "father_code_2",
       "vid": "12345"
    },
    "grey_rollout": {
        "val": false,
        "vid": "45678"
    }
}
```

### getFeatureConfigs
获取用户命中的特定feature的变体详情
```
func (t *AbClient) GetFeatureConfigs(featureId, decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### 参数
| 参数        | 描述         |
|:----------|:-----------|
| featureId | feature Id |

#### 返回值
该函数返回用户命中某个feature的变体详情，通常仅能命中一个变体
```
{
   "feature_key":{
        "val" : "prod",
        "vid" : "20006421"
    }
}
```

### getAllFeatureConfigs
获取用户命中的所有feature的变体详情
```
func (t *AbClient) GetAllFeatureConfigs(decisionId string,
	attributes map[string]interface{}) (map[string]map[string]interface{}, error)
```
#### 参数

#### 返回值
该函数返回用户命中所有feature的变体详情，通常命中多个变体
```
{
   "feature_key":{
       "val" : "prod",
       "vid" : "20006421"
    }
    "feature_key_color":{
       "val" : "true",
       "vid" : "20006423"
    }
}
```

>1、含有“WithImpression”字样的接口均会自动上报曝光事件
>
>2、请务必填写trackId字段，否则会导致上报失效
>
### getExperimentVariantNameWithImpression
同接口“getExperimentVariantName”

### getExperimentConfigsWithImpression
同接口“getExperimentConfigs”

### getAllExperimentConfigsWithImpression
同接口“getAllExperimentConfigs”

### getFeatureConfigsWithImpression
同接口“getFeatureConfigs”

## 其他

### UserAbInfoHandler
用户信息处理接口，冻结实验、进组不出组场景下使用
>1. 使用NewClient初始化AbClient时默认使用空实现，不启用“进组不出组”功能
>2. 继承UserAbInfoHandler接口，自行实现持久化存储；使用NewClientWithUserAbInfo初始化AbClient，并传入自行实现的UserAbInfoHandler类，则可启用“进组不出组”功能

使用Redis缓存示例（仅供参考）
```
client.NewClientWithUserAbInfo("appKey", NewRedisUserAbInfoHandler())

type RedisAbInfoHandler struct{}

func (u *RedisAbInfoHandler) Query(decisionID string) string {
	// need to implement it yourself
    return redis.get(decisionID);
}

func (u *RedisAbInfoHandler) CreateOrUpdate(decisionID, experiment2Variant string) bool {
    // need to implement it yourself
	return redis.set(decisionID, experiment2Variant);
}

func (u *RedisAbInfoHandler) NeedPersistData() bool {
    // return true if customize this interface
	return true
}

func NewRedisUserAbInfoHandler() *RedisAbInfoHandler {
	return &RedisAbInfoHandler{}
}
```

### 匿名上报
>获取不到uuid的用户，可以通过填充device_id或者web_id进行事件上报（私有化场景下也支持bddid）
1. NewClient时设置匿名上报配置，第一个参数（true/开启，false/关闭）匿名上报，第二个参数（true/saas，false/私有化）
```
client.NewClient("appKey", config.WithAnonymousConfig(true, true))
```
2. 添加device_id, web_id, bddid到用户属性attributes，trackId固定传入空字符串""
```
trackId = "";
attributes["device_id"] = 1234; int64
attributes["web_id"] = 5678; int64
attributes["bddid"] = "91011"; string
```
3. 请求activate或其他'WithImpression'接口即可匿名上报