export function trueRandom() {
    return crypto.getRandomValues(new Uint32Array(1))[0] / 2 ** 32;
}

export function randomDesignNameGenerator() {
    return "meshery_compose_" + Math.floor(trueRandom() * 100);
}

