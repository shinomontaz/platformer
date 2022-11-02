#version 330 core
// The first line in glsl source code must always start with a version directive as seen above.

// vTexCoords are the texture coordinates, provided by Pixel
in vec2  vTexCoords;

// fragColor is what the final result is and will be rendered to your screen.
out vec4 fragColor;

// uTexBounds is the texture's boundries, provided by Pixel.
uniform vec4 uTexBounds;

// uTexture is the actualy texture we are sampling from, also provided by Pixel.
uniform sampler2D uTexture;

void main() {
	// t is set to the screen coordinate for the current fragment.
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;
	// And finally, we're telling the GPU that this fragment should be the color as sampled from our texture.
	fragColor = texture(uTexture, t);
}