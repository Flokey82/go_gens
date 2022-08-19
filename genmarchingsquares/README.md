# genmarchingsquares

This package provides a simple implementation of the marching squares algorithm for illustrative purposes for now.

Be warned... Right now I am implementing this blindly without looking anything up, so who knows if what I am doing is correct.

## Scope

Create a package that takes a pixel field and outputs a 2d grid of tiles representing the contours of said field.

### Done
* 2d boolean input field
* 2d output grid of tiles encoded as 4 bit values

### TODO
* Scalar input field
  * Threshold value for contour detection
  * Interpolation of output grid based on scalar values

## Reference
* https://en.wikipedia.org/wiki/Marching_squares
