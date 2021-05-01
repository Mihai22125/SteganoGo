package payload

import "errors"

var ErrUnsupportedBitDepth = errors.New("unsupported bit depth")
var ErrPayloadSize = errors.New("Payload size is bigger than available space on the image")
var ErrNotPayload = errors.New("image didn't contain any payload")
