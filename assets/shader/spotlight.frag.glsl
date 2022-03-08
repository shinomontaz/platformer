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

void main()
{
    vec2 uv = vTexCoords / uTexBounds.zw;

    vec2 light = vec2(uLightX,uLightY) / uTexBounds.zw;
    vec3 pixelColor = texture(uTexture, uv).rgb;

   	float distanceToLight = distance(light.xy, vTexCoords.xy);
    float lightIntencive = ( 1.0 - distanceToLight / 200.0 );
    
    fragColor = vec4(pixelColor * lightIntencive, 1.0);
}
