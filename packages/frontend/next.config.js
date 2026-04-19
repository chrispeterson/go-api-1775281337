/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',

  images: {
    domains: [],
    formats: ['image/avif', 'image/webp'],
  },

  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'Content-Security-Policy',
            value: [
              "default-src 'self'",
              // Allow inline styles for Framer Motion and Tailwind
              "style-src 'self' 'unsafe-inline'",
              // Allow inline scripts for Next.js hydration
              "script-src 'self' 'unsafe-inline' 'unsafe-eval'",
              "img-src 'self' data: blob:",
              "font-src 'self'",
              "connect-src 'self'",
              "media-src 'self'",
              "object-src 'none'",
              "base-uri 'self'",
              "form-action 'self'",
              "frame-ancestors 'none'",
              'upgrade-insecure-requests',
            ].join('; '),
          },
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'Strict-Transport-Security',
            value: 'max-age=63072000; includeSubDomains; preload',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'strict-origin-when-cross-origin',
          },
          {
            key: 'Permissions-Policy',
            value: [
              'camera=()',
              'microphone=()',
              'geolocation=()',
              'browsing-topics=()',
              'interest-cohort=()',
            ].join(', '),
          },
          {
            key: 'X-DNS-Prefetch-Control',
            value: 'on',
          },
        ],
      },
    ];
  },
};

module.exports = nextConfig;
