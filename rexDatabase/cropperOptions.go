package rexDatabase

//	"cropper_options": {
//		"view_mode": 0,
//		"drag_mode": "crop",
//		"initial_aspect_ratio": {
//			"width": 1,
//			"height": 1
//		},
//		"aspect_ratio": {
//			"width": 1,
//			"height": 1
//		},
//		"responsive": true,
//		"restore": true,
//		"check_cross_origin": true,
//		"check_orientation": true,
//		"modal": true,
//		"guides": true,
//		"center": true,
//		"highlight": true,
//		"background": true,
//		"auto_crop": true,
//		"auto_crop_area": true,
//		"movable": true,
//		"rotatable": true,
//		"scalable": true,
//		"zoomable": true,
//		"zoom_on_touch": true,
//		"zoom_on_wheel": true,
//		"wheel_zoom_ratio": true,
//		"crop_box_movable": true,
//		"crop_box_resizable": true,
//		"toggle_drag_mode_on_dblclick": true,
//		"min_container_width": 200,
//		"min_container_height": 100,
//		"min_canvas_width": 0,
//		"min_canvas_height": 0,
//		"min_crop_box_width": 0,
//		"min_crop_box_height": 0,
//		"circled": 0
//

type CropperOptionsViewMode int8

const (
	// 0 -> 无限制
	SysCropperOptionsViewModeNoRestrictions CropperOptionsViewMode = iota
	// 1 -> 限制裁剪框的大小，使其不超过画布的大小
	SysCropperOptionsViewModeRestrictCropBox
	// 2 -> 限制画布的最小尺寸，使其能够容纳在容器内。如果画布和容器的比例不同，最小画布尺寸会在其中一个方向上留出额外的空白
	SysCropperOptionsViewModeRestrictCanvas
	// 3 -> 限制画布的最小尺寸，使其能够填满容器。如果画布和容器的比例不同，容器将无法在某个方向上完全容纳画布
	SysCropperOptionsViewModeRestrictCanvasFillContainer
)

type CropperOptionsDragMode string

const (
	// 创建一个新的裁剪框
	SysCropperOptionsDragModeCrop CropperOptionsDragMode = "crop"
	// 移动画布
	SysCropperOptionsDragModeMove CropperOptionsDragMode = "move"
	// 什么也不做
	SysCropperOptionsDragModeNone CropperOptionsDragMode = "none"
)

type CropperOptions struct {
	ViewMode                         CropperOptionsViewMode    `gorm:"column:view_mode;comment:视图模式;type:smallint;default:1" json:"view_mode"`
	DragMode                         CropperOptionsDragMode    `gorm:"column:drag_mode;comment:拖动模式;type:varchar(12);default:crop" json:"drag_mode"`
	CropperOptionsInitialAspectRatio CropperOptionsAspectRatio `gorm:"embedded;embeddedPrefix:initial_aspect_ratio_" json:"initial_aspect_ratio"`
	CropperOptionsAspectRatio        CropperOptionsAspectRatio `gorm:"embedded;embeddedPrefix:aspect_ratio_" json:"aspect_ratio"`
	//InitialAspectRatioWidth          float64                          `gorm:"column:initial_aspect_ratio_width;comment:初始纵横比宽度;type:float;default:1" json:"initial_aspect_ratio_width"`
	//InitialAspectRatioHeight         float64                          `gorm:"column:initial_aspect_ratio_height;comment:初始纵横比高度;type:float;default:1" json:"initial_aspect_ratio_height"`
	//AspectRatioWidth  float64 `gorm:"column:aspect_ratio_width;comment:纵横比宽度;type:float;default:1" json:"aspect_ratio_width"`
	//AspectRatioHeight float64 `gorm:"column:aspect_ratio_height;comment:纵横比高度;type:float;default:1" json:"aspect_ratio_height"`
	//Data                     []byte                 `gorm:"-" json:"-"`
	//Preview                  []byte                 `gorm:"-" json:"-"`
	Responsive               bool                `gorm:"column:responsive;comment:响应式设计;type:boolean;default:true" json:"responsive"`
	Restore                  bool                `gorm:"column:restore;comment:恢复,调整窗口大小后，恢复被裁剪的区域。;type:boolean;default:true" json:"restore"`
	CheckCrossOrigin         bool                `gorm:"column:check_cross_origin;comment:检查跨域;type:boolean;default:true" json:"check_cross_origin"`
	CheckOrientation         bool                `gorm:"column:check_orientation;comment:检查方向;type:boolean;default:true" json:"check_orientation"`
	Modal                    bool                `gorm:"column:modal;comment:模态,在图像上方、裁剪框下方显示黑色模态框。;type:boolean;default:true" json:"modal"`
	Guides                   bool                `gorm:"column:guides;comment:指南;type:boolean;default:true" json:"guides"`
	Center                   bool                `gorm:"column:center;comment:中心;type:boolean;default:true" json:"center"`
	Highlight                bool                `gorm:"column:highlight;comment:强调，在裁剪框上方显示白色模态框（突出显示裁剪框）;type:boolean;default:true" json:"highlight"`
	Background               bool                `gorm:"column:background;comment:背景,显示容器的网格背景;type:boolean;default:true" json:"background"`
	AutoCrop                 bool                `gorm:"column:auto_crop;comment:自动裁剪;type:boolean;default:true" json:"auto_crop"`
	AutoCropArea             bool                `gorm:"column:auto_crop_area;comment:自动裁剪区域,0-1;type:boolean;default:true" json:"auto_crop_area"`
	Movable                  bool                `gorm:"column:movable;comment:活动,启用此功能可移动图像;type:boolean;default:true" json:"movable"`
	Rotatable                bool                `gorm:"column:rotatable;comment:可旋转;type:boolean;default:true" json:"rotatable"`
	Scalable                 bool                `gorm:"column:scalable;comment:可扩展;type:boolean;default:true" json:"scalable"`
	Zoomable                 bool                `gorm:"column:zoomable;comment:可缩放;type:boolean;default:true" json:"zoomable"`
	ZoomOnTouch              bool                `gorm:"column:zoom_on_touch;comment:触摸缩放;type:boolean;default:true" json:"zoom_on_touch"`
	ZoomOnWheel              bool                `gorm:"column:zoom_on_wheel;comment:车轮上的缩放;type:boolean;default:true" json:"zoom_on_wheel"`
	WheelZoomRatio           bool                `gorm:"column:wheel_zoom_ratio;comment:轮毂缩放比,定义使用鼠标滚轮缩放图像时的缩放比例;type:boolean;default:true" json:"wheel_zoom_ratio"`
	CropBoxMovable           bool                `gorm:"column:crop_box_movable;comment:可移动的裁剪盒,启用此功能后，即可通过拖动来移动裁剪框;type:boolean;default:true" json:"crop_box_movable"`
	CropBoxResizable         bool                `gorm:"column:crop_box_resizable;comment:允许通过拖动来调整裁剪框的大小;type:boolean;default:true" json:"crop_box_resizable"`
	ToggleDragModeOnDblclick bool                `gorm:"column:toggle_drag_mode_on_dblclick;comment:双击切换拖动模式;type:boolean;default:true" json:"toggle_drag_mode_on_dblclick"`
	MinContainerWidth        int64               `gorm:"column:min_container_width;comment:最小容器宽度;type: bigint;default:200" json:"min_container_width"`
	MinContainerHeight       int64               `gorm:"column:min_container_height;comment:最小容器高度;type: bigint;default:100" json:"min_container_height"`
	MinCanvasWidth           int64               `gorm:"column:min_canvas_width;comment:最小canvas宽度;type: bigint;default:0" json:"min_canvas_width"`
	MinCanvasHeight          int64               `gorm:"column:min_canvas_height;comment:最小canvas高度;type: bigint;default:0" json:"min_canvas_height"`
	MinCropBoxWidth          int64               `gorm:"column:min_crop_box_width;comment:最小裁切框宽度;type: bigint;default:0" json:"min_crop_box_width"`
	MinCropBoxHeight         int64               `gorm:"column:min_crop_box_height;comment:最小裁切框高度;type: bigint;default:0" json:"min_crop_box_height"`
	CropperExtraOptions      CropperExtraOptions `gorm:"embedded;embeddedPrefix:extra_" json:"extra"`
}

type CropperOptionsAspectRatio struct {
	Width  float64 `gorm:"column:width;comment:纵横比宽度;type:float;default:1" json:"width"`
	Height float64 `gorm:"column:height;comment:纵横比高度;type:float;default:1" json:"height"`
}

type CropperExtraOptionsCircled int8

const (
	// 圆
	CropperExtraOptionsCircledNo CropperExtraOptionsCircled = iota + 1
	// 非圆
	CropperExtraOptionsCircledYes
)

type CropperExtraOptionsMaxMode int8

const (
	// 圆
	CropperExtraOptionsMaxModeNo CropperExtraOptionsMaxMode = iota + 1
	// 非圆
	CropperExtraOptionsMaxModeYes
)

type CropperExtraOptions struct {
	DefaultImage string                     `gorm:"column:default_image;comment:默认图片;type: varchar(255)" json:"default_image"`
	Circled      CropperExtraOptionsCircled `gorm:"column:circled;comment:是否圆形;type:smallint;default:1" json:"circled"`
	MaxMode      CropperExtraOptionsMaxMode `gorm:"column:max_mode;comment:最大尺寸限制模式;type:smallint;default:1" json:"max_mode"`
	MaxWidth     int64                      `gorm:"column:max_width;comment:最大宽度;type: bigint;default:0" json:"max_width"`
	MaxHeight    int64                      `gorm:"column:max_height;comment:最大高度;type: bigint;default:0" json:"max_height"`
	MinWidth     int64                      `gorm:"column:min_width;comment:最小宽度;type: bigint;default:0" json:"min_width"`
	MinHeight    int64                      `gorm:"column:min_height;comment:最小高度;type: bigint;default:0" json:"min_height"`
	BorderRadius float64                    `gorm:"column:border_radius;comment:边框半径,百分比，保留2位小数，和是否圆形参数冲突;type:float;default:0" json:"border_radius"`
}
