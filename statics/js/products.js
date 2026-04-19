export async function loadProducts(type, subcategory) {
    const res = await fetch(
        `/api/get_products?type=${type}&category=${subcategory}&limit=20&offset=0`
    );

    const data = await res.json();

    const container = document.getElementById("products-container");

    container.innerHTML = data.map(p => `
        <div class="product_card" data-slug="${p.Url}">
            <a href="/product/${p.Url}">
                <img src="/${p.MainPhoto}" class="product_image">
                <p class="product_name_p">${p.Name}</p>

                <div class="product_real_price">
                    ${p.Price} ₽
                </div>

                <div class="product_price">Подробнее</div>
            </a>
        </div>
    `).join("");
}