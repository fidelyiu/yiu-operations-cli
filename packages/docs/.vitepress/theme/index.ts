import { useData, useRoute } from "vitepress";
import DefaultTheme from "vitepress/theme";
import mermaid from "mermaid";
import { nextTick, onMounted, watch } from "vue";
import "./custom.scss";

const renderMermaid = async (isDark: boolean) => {
  if (typeof document === "undefined") {
    return;
  }

  await nextTick();

  const nodes = Array.from(
    document.querySelectorAll<HTMLElement>(".vp-doc .mermaid"),
  );

  if (nodes.length === 0) {
    return;
  }

  mermaid.initialize({
    startOnLoad: false,
    securityLevel: "loose",
    theme: isDark ? "dark" : "default",
  });

  for (const node of nodes) {
    node.removeAttribute("data-processed");
  }

  await mermaid.run({ nodes });
};

export default {
  extends: DefaultTheme,
  setup() {
    const route = useRoute();
    const { isDark } = useData();

    const updateMermaid = async () => {
      await renderMermaid(isDark.value);
    };

    onMounted(() => {
      void updateMermaid();
    });

    watch([() => route.path, () => isDark.value], () => {
      void updateMermaid();
    });
  },
};
