SELECT 
    o.id AS order_id,
    o.created_at,
    u.name AS customer_name,
	u.phone as phone,
    o.total_price,
    o.status
FROM orders o
JOIN users u ON o.user_id = u.id
JOIN order_items oi ON o.id = oi.order_id
JOIN product_variants pv ON oi.variant_id = pv.id
JOIN products p ON pv.product_id = p.id 
GROUP BY o.created_at
ORDER BY o.created_at DESC;


SELECT 
    p.name, 
    pv.value, 
    pv.unit, 
    oi.quantity, 
    oi.price_at_purchase
FROM order_items oi
JOIN product_variants pv ON oi.variant_id = pv.id
JOIN products p ON pv.product_id = p.id
WHERE oi.order_id = 5
