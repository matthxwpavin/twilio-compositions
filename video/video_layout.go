package video

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type VideoLayout struct {
	Resolution  string
	resolutionW uint16
	resolutionH uint16
	regions     []*Region
}

const maxResolutionArea = 921600

// In normal, resolution is not video layout's properties.
// But require for creating video layout will not be invalid restriction
// that relate with composition or composition hooks creating.
// This guarantee creating request will not fail because restriction about video layout.
// Supported solution see compositionhooks package.
func NewVideoLayout(resolution string) (*VideoLayout, error) {
	errRes := errors.New("Error, Invalid resolution.")
	resStr := resolution
	sep := strings.Split(resStr, "x")
	if len(sep) != 2 {
		return nil, errRes
	}

	width, err := strconv.ParseUint(sep[0], 10, 64)
	if err != nil {
		return nil, errRes
	}
	height, err := strconv.ParseUint(sep[1], 10, 64)
	if err != nil {
		return nil, errRes
	}
	if width < XPosDefaultInPix || width > 1280 {
		fmt.Println(width)
		return nil, errRes
	}
	if height < YPosDefaultInPix || height > 1280 {
		fmt.Println(height)
		return nil, errRes
	}

	if area := width * height; area > maxResolutionArea {
		fmt.Println(area)
		return nil, errRes
	}
	return &VideoLayout{
		Resolution:  resolution,
		resolutionW: uint16(width),
		resolutionH: uint16(height),
	}, nil
}

// Usually new video layout from get composition response.
// For access properties with concrete type or via method.
// Getting resolution on video layout from NewFrom is "".
// Because resolution is not video layout's properties.
// More info, see NewVideoLayout.
func NewFrom(obj map[string]interface{}) (*VideoLayout, error) {
	layout := &VideoLayout{}
	for regName, prop := range obj {
		propMap := prop.(map[string]interface{})
		region := &Region{
			Name: regName,
			Prop: &RegionProp{},
		}
		setProps(region.Prop, propMap)
		layout.AddRegion(region)
	}

	return layout, nil
}

func setProps(rgP *RegionProp, pObj map[string]interface{}) error {
	propV := reflect.ValueOf(rgP).Elem()
	propT := propV.Type()
	numF := propT.NumField()
	for k, v := range pObj {
		for i := 0; i < numF; i++ {
			ft := propT.Field(i)
			if k == ft.Tag.Get("alias") {
				fv := propV.Field(i)
				vVal := reflect.ValueOf(v)
				if fv.CanSet() && vVal.IsValid() {
					switch fv.Type().Kind() {
					case reflect.Ptr:
						{
							switch fv.Type().Elem().Kind() {
							case reflect.String:
								nVal := vVal.Interface().(string)
								fv.Set(reflect.ValueOf(&nVal))
							case reflect.Uint16:
								mVal := vVal.Interface().(float64)
								nVal := uint16(mVal)
								fv.Set(reflect.ValueOf(&nVal))
							case reflect.Uint32:
								mVal := vVal.Interface().(float64)
								nVal := uint32(mVal)
								fv.Set(reflect.ValueOf(&nVal))
							case reflect.Int16:
								mVal := vVal.Interface().(float64)
								nVal := int16(mVal)
								fv.Set(reflect.ValueOf(&nVal))
							default:
								return errors.New("Error, pointer to type not expected")
							}
						}
					case reflect.Slice:
						{
							switch fv.Type().Elem().Kind() {
							case reflect.String:
								toSet := []string{}
								nVal := vVal.Interface().([]interface{})
								for _, v := range nVal {
									toSet = append(toSet, v.(string))
								}
								fv.Set(reflect.ValueOf(toSet))
							case reflect.Uint32:
								toSet := []uint32{}
								nVal := vVal.Interface().([]interface{})
								for _, v := range nVal {
									toSet = append(toSet, v.(uint32))
								}
								fv.Set(reflect.ValueOf(toSet))
							default:
								return errors.New("Error, slice of type not expected")
							}
						}
					}
				}
			}
		}
	}
	return nil
}

type Region struct {
	Name string
	Prop *RegionProp
}

func (l *VideoLayout) GetResolution() string {
	return l.Resolution
}

func (l *VideoLayout) AddRegion(reg *Region, regs ...*Region) error {
	return l.addRegion(reg, false, regs...)
}

func (l *VideoLayout) addRegion(reg *Region, validateOn bool, regs ...*Region) error {
	err := errors.New("Error, region must not be nil.")
	if reg == nil {
		return err
	}
	if validateOn {
		if err := l.regionValidation(reg); err != nil {
			return err
		}
	}

	l.regions = append(l.regions, reg)
	if len(regs) != 0 {
		for _, r := range regs {
			if r == nil {
				return err
			}
			if validateOn {
				if err := l.regionValidation(r); err != nil {
					return err
				}
			}

			l.regions = append(l.regions, r)
		}
	}
	return nil
}

const (
	zPosLowerRange = -99
	zPosUpperRange = 99
)

func (l *VideoLayout) regionValidation(reg *Region) error {
	if reg.Name == "" {
		return errors.New("Error, Region's name must not be nil.")
	}

	if reg.Prop.ZPos != nil && (*reg.Prop.ZPos > 0 || *reg.Prop.ZPos < 0) {
		valid, err := l.inRangeValidation(*reg.Prop.ZPos, zPosLowerRange, zPosUpperRange)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("Error, Region's z_pos is invalid.")
		}
	}

	var xPos uint16
	if reg.Prop.Width != nil {
		if reg.Prop.XPos != nil {
			xPos = *reg.Prop.XPos
		}
		valid, err := l.inRangeValidation(*reg.Prop.Width, XPosDefaultInPix, int32(l.resolutionW-xPos))
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("Error, Region's width is invalid.")
		}
	}

	var yPos uint16
	if reg.Prop.Height != nil {
		if reg.Prop.YPos != nil {
			yPos = *reg.Prop.YPos
		}
		valid, err := l.inRangeValidation(*reg.Prop.Height, YPosDefaultInPix, int32(l.resolutionH-yPos))
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("Error, Region's height is invalid.")
		}
	}

	if reg.Prop.MaxColumns != nil {
		valid, err := l.inRangeValidation(*reg.Prop.MaxColumns, 1, 1000)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("Error, Region's max_columns invalid")
		}
	}

	if reg.Prop.MaxRows != nil {
		valid, err := l.inRangeValidation(*reg.Prop.MaxRows, 1, 1000)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("Error, Region's max_rows invalid")
		}
	}

	for _, v := range reg.Prop.CellsExcluded {
		valid, err := l.inRangeValidation(v, 0, 999999)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("Error, Region's cells_excluded invalid")
		}
	}

	if reg.Prop.Reuse != nil {
		switch *reg.Prop.Reuse {
		case ReuseShowOldest, ReuseShowNewest, ReuseNone:
		default:
			return errors.New("Error, Region's reuse invalid")
		}
	}

	if reg.Prop.VideoSources == nil || len(reg.Prop.VideoSources) == 0 {
		return errors.New("Error, Region's video source must specific.")
	}

	return nil
}

func (l *VideoLayout) GetRegions() []*Region {
	return l.regions
}

func (l *VideoLayout) inRangeValidation(value interface{}, lower, upper int32) (bool, error) {
	switch value.(type) {
	case uint16:
		if int32(value.(uint16)) < lower || int32(value.(uint16)) > upper {
			return false, nil
		}
	case int16:
		if int32(value.(int16)) < lower || int32(value.(int16)) > upper {
			return false, nil
		}
	case uint32:
		if int32(value.(uint32)) < lower || int32(value.(uint32)) > upper {
			return false, nil
		}
	default:
		return false, errors.New("Error, value miss match type.")
	}

	return true, nil
}

const (
	XPosDefaultInPix = 16
	YPosDefaultInPix = 16

	ReuseShowOldest = "show_oldest"
	ReuseNone       = "none"
	ReuseShowNewest = "show_newest"
)

type RegionProp struct {

	// X axis value (in pixels) of the region's upper left corner
	// relative to the upper left corner of the Composition viewport.
	// Regions cannot overflow the Composition's area,
	// so x_pos has to be a positive integer less than or equal to
	// the difference between the Composition's width and the width of the region.
	// If the region’s width is missing from the request,
	// it defaults to 16 pixels for this validation.
	XPos *uint16 `json:"x_pos,omitempty" alias:"x_pos"`

	// Y axis value (in pixels) of the region's upper left corner relative to the upper left corner
	// of the Composition viewport. Regions cannot overflow the composition's area,
	// so y_pos has to be a positive integer less than or equal to the difference between
	// the Composition's height and the height of this region. If the region’s height is missing from the request,
	// it defaults to 16 pixels for this validation.
	YPos *uint16 `json:"y_pos,omitempty" alias:"y_pos"`

	// Z position controlling the region's visibility in case of overlaps.
	// Regions with higher values are stacked on top of regions with lower value \
	// for visibility purposes. z_pos must be in the range [-99, 99].
	ZPos *int16 `json:"z_pos,omitempty" alias:"z_pos"`

	// Region's Width. It must be in the range [16, Composition's width - x_pos].
	// This constraint guarantees that the region fits into the Composition's viewport.
	Width *uint16 `json:"width,omitempty" alias:"width"`

	// Region's Height. It must be in the range [16, Composition's height - y_pos].
	// This constraint guarantees that the region fits into the Composition's viewport.
	Height *uint16 `json:"height,omitempty" alias:"height"`

	// Maximum number of columns of the region's placement grid. By default,
	// the region has as many columns as needed to layout all the specified video sources.
	// max_columns must be in the range [1, 1000].
	MaxColumns *uint16 `json:"max_columns,omitempty" alias:"max_columns"`

	// Maximum number of rows of the region's placement grid.
	// By default, the region has as many rows as needed to layout
	// all the specified video sources. max_rows must be in the range [1, 1000].
	MaxRows *uint16 `json:"max_rows,omitempty" alias:"max_rows"`

	// A list of cell indices on the regions layout grid where no video sources can be assigned.
	// Index of first cell (upper left) is 0. Indices grow from left to right and from top to bottom.
	// These values must be in the range [0, 999999].
	CellsExcluded []uint32 `json:"cells_excluded,omitempty" alias:"cells_excluded"`

	// Defines how the region's grid cells are reused for placement purposes. Possible values are:
	//
	//    none: used cells are never reused.
	//    show_oldest: a cell can only be reused when the video source it contains ends.
	//    show_newest: a cell can be reused even if the video source it contains has not ended.
	Reuse *string `json:"reuse,omitempty" alias:"reuse"`

	// The array of video sources that should be placed in this region. All the specified sources must belong to the same Room. It can include:
	//
	//    Zero or more RecordingTrackSid
	//    Zero or more MediaTrackSid
	//    Zero or more ParticipantSid
	//    Zero or more Track names. These can be specified using wildcards (e.g. student*).
	//    The use of [*] has semantics "all if any" meaning zero or more (i.e. all)
	//    depending on whether the target room had video tracks.
	VideoSources []string `json:"video_sources" alias:"video_sources"`

	// An array of video sources to exclude from this region. This region will attempt to display all sources specified in video_sources except for the ones specified in video_sources_excluded. This parameter may include:
	//
	//    Zero or more RecordingTrackSid
	//    Zero or more MediaTrackSid
	//    Zero or more ParticipantSid
	//    Zero or more Track names. These can be specified using wildcards (e.g. student*).
	VideoSourcesExcluded []string `json:"video_sources_excluded,omitempty" alias:"video_sources_excluded"`
}
