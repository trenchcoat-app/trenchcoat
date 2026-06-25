import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import type { Resource, ResourceLanguage } from "i18next";

/**
 * Builds the i18next resources object from Vite's eagerly-imported glob modules.
 *
 * Expects the following directory structure:
 * src/
 * └── i18n
 *     ├── en
 *     │   ├── common.json
 *     │   └── signup.json
 *     └── fr
 *         ├── common.json
 *         └── signup.json
 */
function buildResources(modules: Record<string, { default: ResourceLanguage }>): Resource {
    const resources: Resource = {};

    for (const path in modules) {
        const parts = path.split("/");
        const lng = parts[parts.length - 2];
        const ns = parts[parts.length - 1].replace(".json", "");

        resources[lng] ??= {};
        resources[lng][ns] = modules[path].default;
    }

    return resources;
}

i18n.use(initReactI18next).init({
    lng: "en",
    fallbackLng: "en",
    resources: buildResources(import.meta.glob("../i18n/**/*.json", { eager: true })),
    defaultNS: "common",
    interpolation: { escapeValue: false },
});

export default i18n;
