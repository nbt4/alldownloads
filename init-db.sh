#!/bin/bash
set -e

echo "Initializing AllDownloads database..."

# The database is already available at this point in the init script

# Create tables manually since we can't easily install migrate in this context
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    CREATE TABLE IF NOT EXISTS products (
        id VARCHAR(255) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        vendor VARCHAR(255) NOT NULL,
        category VARCHAR(50) NOT NULL,
        description TEXT,
        icon_url TEXT,
        website_url TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS product_versions (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        product_id VARCHAR(255) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
        version VARCHAR(255) NOT NULL,
        platform VARCHAR(50) NOT NULL,
        architecture VARCHAR(50) NOT NULL,
        download_url TEXT NOT NULL,
        checksum VARCHAR(255),
        checksum_type VARCHAR(50),
        file_size BIGINT,
        filename VARCHAR(255),
        is_latest BOOLEAN DEFAULT FALSE,
        etag VARCHAR(255),
        last_fetched TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        UNIQUE(product_id, version, platform, architecture)
    );

    CREATE TABLE IF NOT EXISTS fetch_jobs (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        product_id VARCHAR(255) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
        status VARCHAR(50) NOT NULL DEFAULT 'pending',
        started_at TIMESTAMP WITH TIME ZONE,
        completed_at TIMESTAMP WITH TIME ZONE,
        error TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_products_vendor ON products(vendor);
    CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
    CREATE INDEX IF NOT EXISTS idx_product_versions_product_id ON product_versions(product_id);
    CREATE INDEX IF NOT EXISTS idx_product_versions_platform ON product_versions(platform);
    CREATE INDEX IF NOT EXISTS idx_product_versions_is_latest ON product_versions(is_latest);
    CREATE INDEX IF NOT EXISTS idx_fetch_jobs_product_id ON fetch_jobs(product_id);
    CREATE INDEX IF NOT EXISTS idx_fetch_jobs_status ON fetch_jobs(status);
    CREATE INDEX IF NOT EXISTS idx_fetch_jobs_created_at ON fetch_jobs(created_at);

    -- Seed data
    INSERT INTO products (id, name, vendor, category, description, icon_url, website_url) VALUES
    ('ubuntu', 'Ubuntu', 'Canonical', 'os', 'Popular Linux distribution known for its ease of use and strong community support', 'https://assets.ubuntu.com/v1/29985a98-ubuntu-logo32.png', 'https://ubuntu.com/'),
    ('debian', 'Debian', 'Debian Project', 'os', 'Stable and secure Linux distribution, the foundation for many other distributions', 'https://www.debian.org/logos/openlogo-nd.svg', 'https://www.debian.org/'),
    ('arch', 'Arch Linux', 'Arch Linux', 'os', 'Lightweight and flexible Linux distribution that follows the KISS principle', 'https://archlinux.org/static/logos/archlinux-logo-dark-scalable.518881f04ca9.svg', 'https://archlinux.org/'),
    ('kali', 'Kali Linux', 'Offensive Security', 'os', 'Penetration testing and security auditing Linux distribution', 'https://www.kali.org/images/kali-logo.svg', 'https://www.kali.org/'),
    ('windows', 'Windows 11', 'Microsoft', 'os', 'Latest version of Microsoft Windows operating system', 'https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE4nqTh', 'https://www.microsoft.com/windows/'),
    ('chrome', 'Google Chrome', 'Google', 'app', 'Fast and secure web browser developed by Google', 'https://www.google.com/chrome/static/images/chrome-logo.svg', 'https://www.google.com/chrome/'),
    ('firefox', 'Mozilla Firefox', 'Mozilla', 'app', 'Open-source web browser focused on privacy and customization', 'https://www.mozilla.org/media/protocol/img/logos/firefox/browser/logo.eb1324e44442.svg', 'https://www.mozilla.org/firefox/'),
    ('brave', 'Brave Browser', 'Brave Software', 'app', 'Privacy-focused web browser that blocks ads and trackers by default', 'https://brave.com/static-assets/images/brave-logo.svg', 'https://brave.com/'),
    ('vscode', 'Visual Studio Code', 'Microsoft', 'tool', 'Lightweight and powerful source code editor with extensive extension support', 'https://code.visualstudio.com/assets/images/code-stable.png', 'https://code.visualstudio.com/'),
    ('termius', 'Termius', 'Termius Corporation', 'tool', 'Cross-platform SSH client with synchronization capabilities', 'https://termius.com/static/uploads/2020/06/icon-512.png', 'https://termius.com/'),
    ('telegram', 'Telegram Desktop', 'Telegram FZ-LLC', 'app', 'Cloud-based messaging app focused on speed and security', 'https://telegram.org/img/t_logo.png', 'https://desktop.telegram.org/'),
    ('whatsapp', 'WhatsApp Desktop', 'Meta', 'app', 'Desktop client for WhatsApp messaging service', 'https://static.whatsapp.net/rsrc.php/v3/yP/r/rYZqPCBaG70.png', 'https://www.whatsapp.com/download'),
    ('tailscale', 'Tailscale', 'Tailscale Inc.', 'tool', 'VPN service that makes devices and applications accessible anywhere in the world', 'https://tailscale.com/kb/1017/install/tailscale-icon.png', 'https://tailscale.com/'),
    ('nextcloud', 'Nextcloud Desktop', 'Nextcloud GmbH', 'tool', 'Desktop sync client for Nextcloud file hosting service', 'https://nextcloud.com/wp-content/uploads/2022/04/nextcloud-logo-blue.svg', 'https://nextcloud.com/'),
    ('notepadplusplus', 'Notepad++', 'Don Ho', 'tool', 'Free source code editor and Notepad replacement for Windows', 'https://notepad-plus-plus.org/images/logo.svg', 'https://notepad-plus-plus.org/'),
    ('powershell', 'PowerShell', 'Microsoft', 'tool', 'Cross-platform task automation solution made up of a command-line shell, scripting language, and configuration management framework', 'https://raw.githubusercontent.com/PowerShell/PowerShell/master/assets/ps_black_64.svg', 'https://github.com/PowerShell/PowerShell')
    ON CONFLICT (id) DO NOTHING;

EOSQL

echo "Database initialization completed successfully."