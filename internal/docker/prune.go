package docker

import (
	"fmt"

	"github.com/moby/moby/client"
)

func (pr *PruneReport) handlePrune(res any, err error) error {
	if err != nil {
		return err
	}

	switch v := res.(type) {
	case client.ContainerPruneResult:
		pr.ContainersPruned += len(v.Report.ContainersDeleted)
		pr.TotalSpaceReclaimed += v.Report.SpaceReclaimed
	case client.ImagePruneResult:
		pr.ImagesPruned += len(v.Report.ImagesDeleted)
		pr.TotalSpaceReclaimed += v.Report.SpaceReclaimed
	case client.VolumePruneResult:
		pr.VolumesPruned += len(v.Report.VolumesDeleted)
		pr.TotalSpaceReclaimed += v.Report.SpaceReclaimed
	case client.NetworkPruneResult:
		pr.NetworksPruned += len(v.Report.NetworksDeleted)
	default:
		return fmt.Errorf("unsupported prune report type: %T", res)
	}

	return nil
}

func (pr PruneReport) Print() {
	fmt.Printf("Prune Report: %+v\n", pr)
}

func (d *Docker) Prune() (PruneReport, error) {
	var report PruneReport

	if err := report.handlePrune(d.client.ContainerPrune(d.ctx, client.ContainerPruneOptions{})); err != nil {
		return report, err
	}

	if err := report.handlePrune(d.client.ImagePrune(d.ctx, client.ImagePruneOptions{})); err != nil {
		return report, err
	}

	if err := report.handlePrune(d.client.VolumePrune(d.ctx, client.VolumePruneOptions{})); err != nil {
		return report, err
	}

	if err := report.handlePrune(d.client.NetworkPrune(d.ctx, client.NetworkPruneOptions{})); err != nil {
		return report, err
	}

	return report, nil

}
