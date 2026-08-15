let lastListUrl = "/";

export function setLastListUrl(url) {
    lastListUrl = url;
}

export function getLastListUrl() {
    return lastListUrl;
}

export function goToList() {
    history.pushState({}, "", lastListUrl);
    window.dispatchEvent(new PopStateEvent("popstate"));
}

export function goToProduct(slug) {
    history.pushState(
        { type: "product", slug },
        "",
        `/product/${slug}`
    );
}

export function getRouteType() {
    return window.location.pathname.startsWith("/product/")
        ? "product"
        : "catalog";
}

export function restoreRoute(loadProductsFn) {
    const path = window.location.pathname;
    const parts = path.split("/");

    // product
    if (path.startsWith("/product/")) return;

    const type = parts[1] || "parfume";
    const sub = parts[2] || "";

    loadProductsFn(type, sub, false);
}