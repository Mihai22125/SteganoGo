package payload

import "errors"

var ErrUnsupportedBitDepth = errors.New("Unsupported bit depth")
var ErrPayloadSize = errors.New("Payload size is bigger than available space on the image")
