import { loader } from "fumadocs-core/source";
import { lucideIconsPlugin } from "fumadocs-core/source/lucide-icons";
import { docs } from "fumadocs-mdx:collections/server";

export const source = loader({
  source: docs.toFumadocsSource(),
  baseUrl: (import.meta.env.PROD ? "/mantrae" : "") + "/docs",
  plugins: [lucideIconsPlugin()],
});
