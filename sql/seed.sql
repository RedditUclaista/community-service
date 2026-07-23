INSERT INTO "COMMUNITY" (id, name, description, rules, banner_url, profile_url, created_by, active, created_at) VALUES
('019f8800-0001-7000-8000-000000000001', 'Informática UCLA', 'Comunidad oficial de estudiantes de Ingeniería de Informática UCLA DCyT.', '1. Respeto mutuo\n2. Compartir recursos académicos', '', '', '019f8726-8be9-7ae5-900e-23f89abae7e9', true, NOW()),
('019f8800-0002-7000-8000-000000000002', 'Desarrollo Web & Móvil', 'Espacio para discutir Flutter, React, Go, Node.js y desarrollo de software.', '1. Código limpio\n2. Ayudar a principiantes', '', '', '019f8726-8be9-7ae5-900e-23f89abae7e9', true, NOW()),
('019f8800-0003-7000-8000-000000000003', 'Inteligencia Artificial & Data', 'Grupo de estudio sobre IA, Machine Learning, Python y Modelos de Lenguaje.', '1. Investigaciones y proyectos', '', '', '019f8726-8be9-7ae5-900e-23f89abae7e9', true, NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO "COMMUNITY_MEMBER" (community_id, user_id, role, joined_at) VALUES
('019f8800-0001-7000-8000-000000000001', '019f8726-8be9-7ae5-900e-23f89abae7e9', 'MODERATOR', NOW()),
('019f8800-0002-7000-8000-000000000002', '019f8726-8be9-7ae5-900e-23f89abae7e9', 'MEMBER', NOW()),
('019f8800-0003-7000-8000-000000000003', '019f8726-8be9-7ae5-900e-23f89abae7e9', 'MEMBER', NOW()),
('4537150b-7868-4013-af98-0307aecac052', '019f8726-8be9-7ae5-900e-23f89abae7e9', 'MODERATOR', NOW())
ON CONFLICT DO NOTHING;
