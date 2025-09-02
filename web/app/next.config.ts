import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "export",
  trailingSlash: true,
  experimental: {
    optimizePackageImports: ["@chakra-ui/react"]
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
