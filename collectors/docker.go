package collectors

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"io"
//	"log"
//	"net/http"
//	"os"
//	"time"
//
//	"github.com/moby/moby/api/types"
//	"github.com/moby/moby/client"
//)
//
//type DockerClient struct {
//	UseSocket bool
//	APIURL    string
//	HTTP      *http.Client
//	CLI       *client.Client
//	Ctx       context.Context
//}
//
//func NewDockerClient(useSocket bool, apiURL string) (*DockerClient, error) {
//	ctx := context.Background()
//	var cli *client.Client
//	var err error
//
//	if useSocket {
//		cli, err = client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
//	} else {
//		cli, err = client.NewClientWithOpts(client.WithHost(apiURL), client.WithAPIVersionNegotiation())
//	}
//
//	if err != nil {
//		return nil, err
//	}
//
//	return &DockerClient{
//		UseSocket: useSocket,
//		APIURL:    apiURL,
//		HTTP:      &http.Client{Timeout: 10 * time.Second},
//		CLI:       cli,
//		Ctx:       ctx,
//	}, nil
//}
//
//func (dc *DockerClient) ListContainers(all bool) ([]types.Container, error) {
//	options := types.ContainerListOptions{All: all}
//	return dc.CLI.ContainerList(dc.Ctx, options)
//}
//
//func (dc *DockerClient) ContainerStats(containerID string) (types.StatsJSON, error) {
//	statsResp, err := dc.CLI.ContainerStats(dc.Ctx, containerID, false)
//	if err != nil {
//		return types.StatsJSON{}, err
//	}
//	defer statsResp.Body.Close()
//
//	var stats types.StatsJSON
//	if err := json.NewDecoder(statsResp.Body).Decode(&stats); err != nil {
//		return stats, err
//	}
//	return stats, nil
//}
//
//func (dc *DockerClient) ContainerDetails(containerID string) (types.ContainerJSON, error) {
//	return dc.CLI.ContainerInspect(dc.Ctx, containerID)
//}
//
//func (dc *DockerClient) FindContainerByPort(targetPort string) ([]map[string]any, error) {
//	containers, err := dc.ListContainers(true)
//	if err != nil {
//		return nil, err
//	}
//
//	var matched []map[string]any
//
//	for _, container := range containers {
//		details, err := dc.ContainerDetails(container.ID)
//		if err != nil {
//			log.Printf("Failed to inspect %s: %v", container.ID, err)
//			continue
//		}
//
//		for port, bindings := range details.HostConfig.PortBindings {
//			for _, b := range bindings {
//				if b.HostPort == targetPort {
//					stats, _ := dc.ContainerStats(container.ID)
//					matched = append(matched, map[string]any{
//						"Id":     container.ID,
//						"Image":  container.ImageID,
//						"Name":   details.Name,
//						"Status": details.State.Status,
//						"Health": details.State.Health,
//						"Stats":  stats,
//					})
//				}
//			}
//		}
//	}
//
//	return matched, nil
//}
//
//func (dc *DockerClient) GetContainerHealth(containerName string) map[string]any {
//	container, err := dc.CLI.ContainerInspect(dc.Ctx, containerName)
//	if err != nil {
//		return map[string]any{"error": err.Error()}
//	}
//	return map[string]any{
//		"name":   container.Name,
//		"status": container.State.Status,
//		"health": container.State.Health,
//	}
//}
//
//func (dc *DockerClient) Ping() bool {
//	_, err := dc.CLI.Ping(dc.Ctx)
//	return err == nil
//}
