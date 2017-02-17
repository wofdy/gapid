#version 310 es
precision mediump float;
precision highp int;

layout(binding = 1) uniform mediump texture2D uTexture;
layout(binding = 0) uniform mediump sampler uSampler;
layout(binding = 4) uniform mediump texture2DArray uTextureArray;
layout(binding = 3) uniform mediump textureCube uTextureCube;
layout(binding = 2) uniform mediump texture3D uTexture3D;

layout(location = 0) in vec2 vTex;
layout(location = 1) in vec3 vTex3;
layout(location = 0) out vec4 FragColor;

vec4 sample_func(mediump sampler samp, vec2 uv)
{
    return texture(sampler2D(uTexture, samp), uv);
}

vec4 sample_func_dual(mediump sampler samp, mediump texture2D tex, vec2 uv)
{
    return texture(sampler2D(tex, samp), uv);
}

void main()
{
    vec2 off = (vec2(1.0) / vec2(textureSize(sampler2D(uTexture, uSampler), 0)));
    vec2 off2 = (vec2(1.0) / vec2(textureSize(sampler2D(uTexture, uSampler), 1)));
    highp vec2 param = ((vTex + off) + off2);
    vec4 c0 = sample_func(uSampler, param);
    highp vec2 param_1 = ((vTex + off) + off2);
    vec4 c1 = sample_func_dual(uSampler, uTexture, param_1);
    vec4 c2 = texture(sampler2DArray(uTextureArray, uSampler), vTex3);
    vec4 c3 = texture(samplerCube(uTextureCube, uSampler), vTex3);
    vec4 c4 = texture(sampler3D(uTexture3D, uSampler), vTex3);
    FragColor = ((((c0 + c1) + c2) + c3) + c4);
}

