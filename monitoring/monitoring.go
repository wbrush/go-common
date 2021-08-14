package monitoring

import (
    "cloud.google.com/go/profiler"
    "github.com/sirupsen/logrus"
)

var monitor  *Monitor
type Monitor struct {
    ServiceName string
    Version     string
    ProjectId   string
    Active      bool

    // ?
}

func NewMonitor(service, version, projectId string) (monitor *Monitor) {
    if monitor == nil {
        monitor = &Monitor{
            ServiceName: service,
            Version:     version,
            ProjectId:   projectId,
            Active:      true,
        }
    } else {
        logrus.Errorf("Monitoring already configured!")
        return monitor
    }

    if service == "" {
        logrus.Errorf("Service must not be empty")
        monitor.Active = false
        return nil
    } //  version can be empty

    cfg := profiler.Config{
        Service:        monitor.ServiceName,
        ServiceVersion: monitor.Version,
    }
    if projectId != "" {
        cfg.ProjectID = projectId
    }

    // Profiler initialization, best done as early as possible.
    if err := profiler.Start(cfg); err != nil {
        logrus.Errorf("Cannot start GCP Profiler: %s", err.Error())
        monitor.Active = false
    }

    return monitor
}

func GetMonitor() *Monitor {
    return monitor
}
