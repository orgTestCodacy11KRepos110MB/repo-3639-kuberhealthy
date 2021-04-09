// Copyright 2018 Comcast Cable Communications Management, LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"

	kh "github.com/Comcast/kuberhealthy/v2/pkg/checks/external/checkclient"
	log "github.com/sirupsen/logrus"
)

// parseDebugSettings parses debug settings and fatals on errors.
func parseDebugSettings() {

	// Enable debug logging if required.
	if len(debugEnv) != 0 {
		var err error
		debug, err = strconv.ParseBool(debugEnv)
		if err != nil {
			log.Fatalln("failed to parse DEBUG environment variable:", err.Error())
		}
	}

	// Turn on debug logging.
	if debug {
		log.Infoln("Debug logging enabled.")
		log.SetLevel(log.DebugLevel)
	}
	log.Debugln(os.Args)
}

// parseInputValues parses all incoming environment variables for the program into globals and fatals on errors.
func parseInputValues() {

	// Calculated in binary SI units (75 * 1024^2 = 75Mi memory).
	checkName = defaultCheckName
	if len(checkNameEnv) != 0 {
		checkName = checkNameEnv
		log.Infoln("Parsed CHECK_NAME:", checkName)
	}

	// Parse incoming namespace environment variable
	checkNamespace = defaultCheckNamespace
	data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Warnln("Failed to open namespace file:", err.Error())
	}
	if len(data) != 0 {
		log.Infoln("Found pod namespace:", string(data))
		checkNamespace = string(data)
	}
	if len(checkNamespaceEnv) != 0 {
		checkNamespace = checkNamespaceEnv
		log.Infoln("Parsed CHECK_NAMESPACE:", checkNamespace)
	}
	log.Infoln("Performing check in", checkNamespace, "namespace.")

	// // Parse incpoming deployment tolerations
	// if len(checkDeploymentTolerationsEnv) > 0 {
	// 	splitEnvVars := strings.Split(checkDeploymentTolerationsEnv, ",")
	// 	for _, splitEnvVarKeyValuePair := range splitEnvVars {
	// 		parsedEnvVarKeyValuePair := strings.Split(splitEnvVarKeyValuePair, "=")
	// 		if len(parsedEnvVarKeyValuePair) != 2 {
	// 			log.Warnln("Unable to parse key value pair:", splitEnvVarKeyValuePair)
	// 			log.Warnln("Setting operator to", corev1.TolerationOpExists)
	// 			t := corev1.Toleration{
	// 				Key:      parsedEnvVarKeyValuePair[0],
	// 				Operator: corev1.TolerationOpExists,
	// 			}
	// 			log.Infoln("Adding toleration to deployment:", t)
	// 			checkDeploymentTolerations = append(checkDeploymentTolerations, t)
	// 			continue
	// 		}
	// 		parsedEnvVarValueEffect := strings.Split(parsedEnvVarKeyValuePair[1], ":")
	// 		if len(parsedEnvVarValueEffect) != 2 {
	// 			log.Warnln("Unable to parse complete toleration value and effect:", parsedEnvVarValueEffect)
	// 			t := corev1.Toleration{
	// 				Key:      parsedEnvVarKeyValuePair[0],
	// 				Operator: corev1.TolerationOpEqual,
	// 				Value:    parsedEnvVarKeyValuePair[1],
	// 			}
	// 			log.Infoln("Adding toleration to deployment:", t)
	// 			checkDeploymentTolerations = append(checkDeploymentTolerations, t)
	// 			continue
	// 		}
	// 		t := corev1.Toleration{
	// 			Key:      parsedEnvVarKeyValuePair[0],
	// 			Operator: corev1.TolerationOpEqual,
	// 			Value:    parsedEnvVarValueEffect[0],
	// 			Effect:   corev1.TaintEffect(parsedEnvVarValueEffect[1]),
	// 		}
	// 		log.Infoln("Adding toleration to deployment:", t)
	// 		checkDeploymentTolerations = append(checkDeploymentTolerations, t)
	// 	}
	// 	log.Infoln("Parsed TOLERATIONS:", checkDeploymentTolerations)
	// }

	// Parse incoming deployment node selectors
	if len(checkNodeSelectorsEnv) > 0 {
		// splitEnvVars := strings.Split(checkNodeSelectorsEnv, ",")
		// for _, splitEnvVarKeyValuePair := range splitEnvVars {
		// 	parsedEnvVarKeyValuePair := strings.Split(splitEnvVarKeyValuePair, "=")
		// 	if len(parsedEnvVarKeyValuePair) != 2 {
		// 		log.Warnln("Unable to parse key value pair:", splitEnvVarKeyValuePair)
		// 		continue
		// 	}
		// 	if _, ok := checkDeploymentNodeSelectors[parsedEnvVarKeyValuePair[0]]; !ok {
		// 		checkDeploymentNodeSelectors[parsedEnvVarKeyValuePair[0]] = parsedEnvVarKeyValuePair[1]
		// 	}
		// }
		log.Infoln("Parsed NODE_SELECTOR:", checkNodeSelectorsEnv)
		// log.Infoln("Parsed NODE_SELECTOR:", checkDeploymentNodeSelectors)
	}

	// Parse incoming node selectors
	// if len(checkNodeSelectorsEnv) > 0 {
	// 	splitEnvVars := strings.Split(checkNodeSelectorsEnv, ",")
	// 	for _, splitEnvVarKeyValuePair := range splitEnvVars {

	// 		// Split each comma-separated input based on `=`
	// 		parsedEnvVarKeyValuePair := strings.Split(splitEnvVarKeyValuePair, "=")
	// 		if len(parsedEnvVarKeyValuePair) != 2 {
	// 			log.Warnln("Unable to parse key value pair:", splitEnvVarKeyValuePair)
	// 		}

	// 		if len(parsedEnvVarKeyValuePair) == 2 {
	// 			if _, ok := checkNodeSelectors[parsedEnvVarKeyValuePair[0]]; !ok {
	// 				checkNodeSelectors[parsedEnvVarKeyValuePair[0]] = parsedEnvVarKeyValuePair[1]
	// 			}
	// 		}

	// 		// Split each comma-separated input based on `:`
	// 		parsedEnvVarKeyValuePair = strings.Split(splitEnvVarKeyValuePair, ":")
	// 		if len(parsedEnvVarKeyValuePair) != 2 {
	// 			log.Warnln("Unable to parse key value pair:", splitEnvVarKeyValuePair)
	// 		}

	// 		if len(parsedEnvVarKeyValuePair) == 2 {
	// 			if _, ok := checkNodeSelectors[parsedEnvVarKeyValuePair[0]]; !ok {
	// 				checkNodeSelectors[parsedEnvVarKeyValuePair[0]] = parsedEnvVarKeyValuePair[1]
	// 			}
	// 		}
	// 	}
	// 	log.Infoln("Parsed NODE_SELECTOR:", checkNodeSelectors)
	// }

	// // Parse incoming check pod resource requests and limits
	// // Calculated in decimal SI units (15 = 15m cpu).
	// millicoreRequest = defaultMillicoreRequest
	// if len(millicoreRequestEnv) != 0 {
	// 	cpuRequest, err := strconv.ParseInt(millicoreRequestEnv, 10, 64)
	// 	if err != nil {
	// 		log.Fatalln("error occurred attempting to parse CHECK_POD_CPU_REQUEST:", err)
	// 	}
	// 	millicoreRequest = int(cpuRequest)
	// 	log.Infoln("Parsed CHECK_POD_CPU_REQUEST:", millicoreRequest)
	// }

	// // Calculated in decimal SI units (75 = 75m cpu).
	// millicoreLimit = defaultMillicoreLimit
	// if len(millicoreLimitEnv) != 0 {
	// 	cpuLimit, err := strconv.ParseInt(millicoreLimitEnv, 10, 64)
	// 	if err != nil {
	// 		log.Fatalln("error occurred attempting to parse CHECK_POD_CPU_LIMIT:", err)
	// 	}
	// 	millicoreLimit = int(cpuLimit)
	// 	log.Infoln("Parsed CHECK_POD_CPU_LIMIT:", millicoreLimit)
	// }

	// // Calculated in binary SI units (20 * 1024^2 = 20Mi memory).
	// memoryRequest = defaultMemoryRequest
	// if len(memoryRequestEnv) != 0 {
	// 	memRequest, err := strconv.ParseInt(memoryRequestEnv, 10, 64)
	// 	if err != nil {
	// 		log.Fatalln("error occurred attempting to parse CHECK_POD_MEM_REQUEST:", err)
	// 	}
	// 	memoryRequest = int(memRequest) * 1024 * 1024
	// 	log.Infoln("Parsed CHECK_POD_MEM_REQUEST:", memoryRequest)
	// }

	// // Calculated in binary SI units (75 * 1024^2 = 75Mi memory).
	// memoryLimit = defaultMemoryLimit
	// if len(memoryLimitEnv) != 0 {
	// 	memLimit, err := strconv.ParseInt(memoryLimitEnv, 10, 64)
	// 	if err != nil {
	// 		log.Fatalln("error occurred attempting to parse CHECK_POD_MEM_LIMIT:", err)
	// 	}
	// 	memoryLimit = int(memLimit) * 1024 * 1024
	// 	log.Infoln("Parsed CHECK_POD_MEM_LIMIT:", memoryLimit)
	// }

	// Set check time limit to default
	checkTimeLimit = defaultCheckTimeLimit
	// Get the deadline time in unix from the env var
	timeDeadline, err := kh.GetDeadline()
	if err != nil {
		log.Infoln("There was an issue getting the check deadline:", err.Error())
	}
	checkTimeLimit = timeDeadline.Sub(time.Now().Add(time.Second * 5))
	log.Infoln("Check time limit set to:", checkTimeLimit)

	// Parse incoming custom shutdown grace period seconds
	shutdownGracePeriod = defaultShutdownGracePeriod
	if len(shutdownGracePeriodEnv) != 0 {
		duration, err := time.ParseDuration(shutdownGracePeriodEnv)
		if err != nil {
			log.Fatalln("error occurred attempting to parse SHUTDOWN_GRACE_PERIOD:", err)
		}
		if duration.Seconds() < 1 {
			log.Fatalln("error occurred attempting to parse SHUTDOWN_GRACE_PERIOD.  A value less than 1 was parsed:", duration.Seconds())
		}
		shutdownGracePeriod = duration
		log.Infoln("Parsed SHUTDOWN_GRACE_PERIOD:", shutdownGracePeriod)
	}
}