#version 330 core

// https://www.shadertoy.com/view/7dlXWM

// Whether or not shadows can hide objects
//#define OBSTRUCT

in vec2 vTexCoords;
out vec4 fragColor;

// Pixel default uniforms
uniform vec4      uTexBounds;
uniform sampler2D uTexture;

// Our custom uniforms
uniform float uTime;
uniform vec2 uLight;
//uniform float uObjects[1000];
uniform vec4 uObjects[1000];

uniform int uNumObjects;

// Utilities
#define drawSDF(dist, col) color = mix(color, col, smoothstep(1.0, 0.0, dist))

struct ShadowVol2D {
    vec2 ap;
    vec2 ad;
    vec2 bp;
    vec2 bd;
};

// Shadow volumes
ShadowVol2D shadowVolBox(in vec2 l, in vec2 b) {
    vec2 s = vec2(l.x < 0.0 ? -1.0 : 1.0, l.y < 0.0 ? -1.0 : 1.0);
    vec2 c1 = vec2(b.x * sign(b.y - abs(l.y)), b.y) * s;
    vec2 c2 = vec2(b.x, b.y * sign(b.x - abs(l.x))) * s;
    return ShadowVol2D(c1, normalize(c1 - l), c2, normalize(c2 - l));
}

float sdBox(in vec2 p, in vec2 b) {
    vec2 q = abs(p) - b;
    return length(max(q, 0.0)) + min(0.0, max(q.x, q.y));
}

float sdDisc(in vec2 p, in float r) {
    return length(p) - r;
}

float sdBox2(in vec2 uv, in vec2 tl, in vec2 br) {
    vec2 d = max(tl-uv, uv-br);
    return length(max(vec2(0.0), d)) + min(0.0, max(d.x, d.y));
}

// a, b - points
// ld - left down corner of a rectangle
// ru - right up corner of a rectangle
bool isShadowedByBox( vec2 a, vec2 b, vec2 ld, vec2 ru ) {

    if ( ( a.x > ru.x && b.x > ru.x ) || ( a.x < ld.x && b.x < ld.x ) ) {
        return false;
    }

    if ( ( a.y > ru.y && b.y > ru.y ) || ( a.y < ld.y && b.y < ld.y ) ) {
        return false;
    }

    if ( b.y <= ru.y && b.y >= ld.y && b.x <= ru.x && b.x >= ld.x ) { // this check that point inside box
        return false;
    }

    // get line equation
    float A = b.y - a.y;
    float B = - ( b.x - a.x );
    float C = - a.x * ( b.y - a.y ) + a.y * ( b.x - a.x );

    // check 4 signs to test if all vertexes lies in same halfplane defined by line
    float lds = sign( A * ld.x + B * ld.y + C);
    float rds = sign( A * ru.x + B * ld.y + C);
    float rus = sign( A * ru.x + B * ru.y + C);
    float lus = sign( A * ld.x + B * ru.y + C);

    if (lds == rds && lds == rus && lds == lus ) {
        return false;
    }

    // now check a or b inside a square
    // we test if a or b lies in different halfplanes defined by each rectangle border; here we know line a,b intersects rect somewhere
    float A1 = ru.y - ru.y;
    float B1 = - ( ru.x - ld.x );
    float C1 = - ld.x * ( ru.y - ru.y ) + ru.y * ( ru.x - ld.x );
    float as = sign( A1 * a.x + B1 * a.y + C1);
    float bs = sign( A1 * b.x + B1 * b.y + C1);
    if ( as != bs ) {
        return true;
    }

    float A2 = ld.y - ru.y;
    float B2 = - ( ru.x - ru.x );
    float C2 = - ru.x * ( ld.y - ru.y ) + ru.y * ( ru.x - ru.x );
    as = sign( A2 * a.x + B2 * a.y + C2);
    bs = sign( A2 * b.x + B2 * b.y + C2);
    if ( as != bs ) {
        return true;
    }

    float A3 = ld.y - ld.y;
    float B3 = - ( ru.x - ld.x );
    float C3 = - ld.x * ( ld.y - ld.y ) + ld.y * ( ru.x - ld.x );
    as = sign( A3 * a.x + B3 * a.y + C3);
    bs = sign( A3 * b.x + B3 * b.y + C3);
    if ( as != bs ) {
        return true;
    }

    float A4 = ld.y - ru.y;
    float B4 = - ( ld.x - ld.x );
    float C4 = - ld.x * ( ld.y - ru.y ) + ld.y * ( ld.x - ld.x );
    as = sign( A4 * a.x + B4 * a.y + C4);
    bs = sign( A4 * b.x + B4 * b.y + C4);
    if ( as != bs ) {
        return true;
    }

    return false;
}

bool isInsideBox( vec2 b, vec2 ld, vec2 ru ) {

    if ( ( b.x > ru.x ) || ( b.x < ld.x ) ) {
        return false;
    }

    if ( ( b.y > ru.y ) || ( b.y < ld.y ) ) {
        return false;
    }

    return true;
}

float fillMask(float dist)
{
	return clamp(-dist, 0.0, 1.0);
}

float circleDist(vec2 p, float radius)
{
	return length(p) - radius;
}

vec3 drawLight(vec2 p, vec2 pos, vec3 color, float range, float radius)
{
	// distance to light
	float ld = length(p - pos);
	
	// out of range
	if (ld > range) return vec3(0.0);
	
	// shadow and falloff
	float fall = (range - ld)/range;
	fall *= fall;
	float source = fillMask(circleDist(p - pos, radius));
	return (fall + source) * color;
}


// fragCoord -> vTexCoords
// iResolution.xy -> uTexBounds.zw
// mainImage(out vec4 fragColor, in vec2 fragCoord) -> main() + definition of 
// in vec2 vTexCoords;
// out vec4 fragColor;

void main() {
    vec2 uv = vTexCoords.xy / uTexBounds.zw;
    vec4 pixelColor = texture(uTexture, uv);

    vec2 toLight = uv - uLight/uTexBounds.zw;
    vec3 color = pixelColor.rgb; // / (dot(toLight, toLight) + pixelColor.rgb);

    vec2 circle_pos = uLight + uTexBounds.zw/2;

    bool shadowed = false;
    bool inside = false;
    // Shapes and shadow volumes
    for (int i=0; i<uNumObjects; i++) {
        vec2 bp = vec2(uObjects[i].x, uObjects[i].y) + uTexBounds.zw/2;     // []Vec4
        vec2 bb = vec2(uObjects[i].z, uObjects[i].w) + uTexBounds.zw/2;

       if (isInsideBox( vTexCoords.xy, bp, bb )) {
            inside = true;
            continue;
        }

       if (isShadowedByBox( circle_pos, vTexCoords.xy, bp, bb )) {
            shadowed = true;
        }
    }

    if (!inside) {
        if (!shadowed) {
            color += drawLight(vTexCoords.xy, circle_pos, vec3(1.0, 0.75, 0.5), 100.0, 10.0);
        } else {
            color -= 0.1;
        }

    }

    fragColor = vec4(color, 1.0);
}