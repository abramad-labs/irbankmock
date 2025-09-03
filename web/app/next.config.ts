import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "export",
  trailingSlash: true,
  experimental: {
    optimizePackageImports: ["@chakra-ui/react"],
    caseSensitiveRoutes: false
  },
  async rewrites() {
    return [
      {
        source: "/:path*",
        destination: "http://localhost:3000/:path*",
      }
    ]
  },
};

export default nextConfig;
