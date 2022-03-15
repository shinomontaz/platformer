#version 330 core

// inspired by https://www.shadertoy.com/view/ltKfWz

in vec2 vTexCoords;
out vec4 fragColor;

// Pixel default uniforms
uniform vec4      uTexBounds;
uniform sampler2D uTexture;

// Our custom uniforms
uniform float uLightX;
uniform float uLightY;
uniform float uTime;

// TODO: use time

void main()
{
    vec2 uv = vTexCoords / uTexBounds.zw;

    vec2 light = vec2(uLightX,uLightY);
    vec4 pixelColor = texture(uTexture, uv.xy);

   	float distanceToLight = distance(light.xy, vTexCoords.xy);
    float lightIntencive = ( 1.0 - distanceToLight / 75.0 * ( 2 + sin(uTime) ) );
    
    fragColor = pixelColor * lightIntencive;
}
