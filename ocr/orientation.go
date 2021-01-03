package ocr

import pb "google.golang.org/genproto/googleapis/cloud/vision/v1"

type Orientation int

const (
	OrientationUnknown Orientation = -1
	Orientation0       Orientation = 0
	Orientation90      Orientation = 90
	Orientation270     Orientation = -90
	Orientation180     Orientation = 180
)

func DetectOrientation(in []*pb.EntityAnnotation) Orientation {
	if in == nil {
		return OrientationUnknown
	}

	counts := map[Orientation]int{
		Orientation0:   0,
		Orientation90:  0,
		Orientation270: 0,
		Orientation180: 0,
	}

	for _, entity := range in {
		if entity == nil {
			continue
		}
		boundary := entity.GetBoundingPoly()
		if boundary == nil || len(boundary.GetVertices()) != 4 {
			continue
		}

		centerX, centerY := int32(0), int32(0)
		vertices := boundary.GetVertices()
		for _, point := range vertices {
			centerX += point.X
			centerY += point.Y
		}

		centerX = centerX / 4
		centerY = centerY / 4

		x0 := vertices[0].X
		y0 := vertices[0].Y

		if x0 < centerX {
			if y0 < centerY {

				counts[Orientation0] += 1
			} else {
				counts[Orientation270] += 1
			}
		} else {
			if y0 < centerY {
				counts[Orientation90] += 1
			} else {
				counts[Orientation180] += 1
			}
		}
	}

	return maxCount(counts)
}

func maxCount(counts map[Orientation]int) Orientation {
	max := OrientationUnknown
	for k, v := range counts {
		if v > counts[max] {
			max = k
		}
	}
	return max
}
