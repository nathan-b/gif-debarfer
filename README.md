# GIF debarfer
With more and more programs supporting dark mode, some animated emojis end up looking pretty terrible. They have a "halo" of blurry light pixels surrounding the emoji itself, which is largely invisible against a white background but stands out very unattractively against a dark background.

This program attempts to algorithmically remove these artifacts using a revolutionary process called "debarfing". The debarf algorithm is very simple:
* Classify each pixel as either "dark", "light", or "transparent".
* Find every light pixel between a dark pixel and a transparent pixel and turn it transparent.

Obviously the algorithm is not perfect, and will not replace a skilled editor using a competently powerful graphics editing program. However, for images (both static and animated GIFs) afflicted with the problem described above, it will at least make them more tolerable to behold.
