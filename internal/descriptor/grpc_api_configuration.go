package descriptor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/monime-lab/grpc-gateway/v2/internal/descriptor/apiconfig"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

func loadGrpcAPIServiceFromYAML(yamlFileContents []byte, yamlSourceLogName string) (*apiconfig.GrpcAPIService, error) {
	var yamlContents interface{}
	err := yaml.Unmarshal(yamlFileContents, &yamlContents)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gRPC API Configuration from YAML in '%v': %v", yamlSourceLogName, err)
	}

	jsonContents, err := json.Marshal(yamlContents)
	if err != nil {
		return nil, err
	}

	// As our GrpcAPIService is incomplete, accept unknown fields.
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	serviceConfiguration := apiconfig.GrpcAPIService{}
	if err := unmarshaler.Unmarshal(jsonContents, &serviceConfiguration); err != nil {
		return nil, fmt.Errorf("failed to parse gRPC API Configuration from YAML in '%v': %v", yamlSourceLogName, err)
	}

	return &serviceConfiguration, nil
}

func registerHTTPRulesFromGrpcAPIService(registry *Registry, service *apiconfig.GrpcAPIService, sourceLogName string) error {
	if service.Http == nil {
		// Nothing to do
		return nil
	}

	for _, rule := range service.Http.GetRules() {
		selector := "." + strings.Trim(rule.GetSelector(), " ")
		if strings.ContainsAny(selector, "*, ") {
			return fmt.Errorf("selector '%v' in %v must specify a single service method without wildcards", rule.GetSelector(), sourceLogName)
		}

		registry.AddExternalHTTPRule(selector, rule)
	}

	return nil
}

// LoadGrpcAPIServiceFromYAML loads a gRPC API Configuration from the given YAML file
// and registers the HttpRule descriptions contained in it as externalHTTPRules in
// the given registry. This must be done before loading the proto file.
//
// You can learn more about gRPC API Service descriptions from google's documentation
// at https://cloud.google.com/endpoints/docs/grpc/grpc-service-config
//
// Note that for the purposes of the gateway generator we only consider a subset of all
// available features google supports in their service descriptions.
func (r *Registry) LoadGrpcAPIServiceFromYAML(yamlFile string) error {
	yamlFileContents, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read gRPC API Configuration description from '%v': %v", yamlFile, err)
	}

	service, err := loadGrpcAPIServiceFromYAML(yamlFileContents, yamlFile)
	if err != nil {
		return err
	}

	return registerHTTPRulesFromGrpcAPIService(r, service, yamlFile)
}
