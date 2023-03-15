# lemosaic

Yet another one mosaic constructor made for fun. Can either [-mosaic] with squared images or [-pour] with uniform colours or  get RGBA average size of a picture with [-average].

Size of a squares depends of [-chunk] flag. Default is 20.

Capable of processing images in *webp, jpg, png, tiff, gif* formats and output it in *png, jpg, gif and tiff*.
Source images are specified with [-source] flag. Default is "./pics". All above mentioned formats are acceptable, but resizing srcs to square *will* produce weird results if src resolution is not atleast close to a square.

Default output is jpeg. Default jpeg quality is 50 - settable with [-qual] flag [1-100]. Output format can be changed with [-encoder] flag [0-3] - jpeg, png, tiff, gif.

Overall output quality and internal encoding settable with [-interpol] ([0-3] From worst to best: NearestNeighbor; ApproxBiLinear; CatmullRom; BiLinear;) and with [-format] flag ([0-5] : RGBA; RGBA64; NRGBA; NRGBA64; GRAY; CMYK;). Defaults are 0 in both cases.

To ensure an actual *squaring* of an image with flag [-normal] size of outputted image will be adjusted to actual *(actual <= original)*. Without this flag on the edges will be cutted squares.

Semi-controlable number of goroutines in use with flag [-routine]. It will be internally adjusted to an actual number *(actual <= requested)* to speed up the process and to ensure that routines do not overlap on target pixels-slice with size of a square in mind. Default is 500, although most of the time no more than 40 will be actually launched. The idea behind that based on simple image resolution progression of (n\*2) pixels and the need to flatten this to (2\*n) decoder/encoder bottleneck. #todo - number of routines is too low, calculations are wrong

As the goal was to accept and work with multiple types, optimisations for specific formats are ignored. May be easily done with unsafe slice operations and loop unrolling, but I wanted to work with at most *.Drawable* interface without specifying underlying format.

Default img size limit is ~10mb. Use [-unmax] to unset internal limit and rely on "image" lib to decide what's too big and what's not.

As it's just silly little project, default behavior is local mosaic - accepts an image [-mosaic|-pour] [target] and outputs it in ./target/ dir. With [-serve] flag it will instead run on localhost:8080. HTML in use is from Chang S. - Go-Web-Programming.

<img src="https://i.imgur.com/CNDJscB.png" width="250">
