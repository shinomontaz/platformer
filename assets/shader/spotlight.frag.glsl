#version 330 core

#define LIGHT_RANGE 90.

in vec2 vTexCoords;
out vec4 fragColor;

// Pixel default uniforms
uniform vec4      uTexBounds;
uniform sampler2D uTexture;

// Our custom uniforms
uniform float uTime;

void main()
{
    vec2 uv = gl_FragCoord.xy / uTexBounds.zw;

    vec2 light = vec2(.2,.2);
    
    light = vec2(abs(sin(uTime)),.2);
    
    vec3 finalColor = vec3(.8,.8,.8) * pow(max(dot(normalize(light),normalize(uv)),0.),LIGHT_RANGE);
    vec3 bg = texture(uTexture, uv).xyz / 4.;
    
	fragColor = vec4(bg + finalColor.xyz,1.);
}