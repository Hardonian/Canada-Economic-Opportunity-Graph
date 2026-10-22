import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        background: "#0B1520",     // Deep Atlantic navy
        surface: "#122334",        // High-contrast slate surface
        card: "#182D3B",           // Institutional blue-slate panel
        cardHover: "#1F3A48",      // Clear interactive hover state
        border: "#547082",         // 3:1+ control boundary on supported panels
        borderSubtle: "#304D60",   // Decorative divider, not a sole status cue
        primary: {
          DEFAULT: "#55D6BE",      // Atlantic teal
          hover: "#75E3D0",
          light: "#A2F0E2",
        },
        aurora: {
          DEFAULT: "#55D6BE",      // Atlantic teal
          glow: "#55D6BE",
          emerald: "#38B99F",
          teal: "#55D6BE",
          mint: "#A2F0E2",
        },
        accent: {
          cyan: "#55D6BE",         // Atlantic teal
          gold: "#E7B96B",         // Canadian mineral gold
          green: "#55D6BE",        // Atlantic teal
          purple: "#C5A3FF",       // Northern Light Violet
          copper: "#D58D62",       // Canadian copper
          crimson: "#E56B5D",      // Sovereign maple red
        },
        gold: {
          DEFAULT: "#E7B96B",      // Canadian mineral gold
          light: "#F5D99D",
          dark: "#C58F45",
          copper: "#D58D62",
        },
        crimson: {
          DEFAULT: "#E56B5D",      // Canadian maple accent; use sparingly
          light: "#F1A096",
          dark: "#B94E45",
        },
        text: {
          main: "#F7F4ED",         // Warm quartz
          muted: "#C3D0D2",        // 8:1+ across supported panels
          subtle: "#9BB2B9",       // 5.5:1+ across supported panels
        }
      },
      fontFamily: {
        sans: ["Inter", "-apple-system", "BlinkMacSystemFont", "Segoe UI", "sans-serif"],
        mono: ["JetBrains Mono", "SFMono-Regular", "Menlo", "monospace"],
      },
      backgroundImage: {
        "aurora-mesh": "radial-gradient(ellipse 80% 50% at 50% -20%, rgba(85, 214, 190, 0.16), rgba(11, 21, 32, 0))",
        "gold-glow": "radial-gradient(ellipse 60% 40% at 50% 120%, rgba(231, 185, 107, 0.1), rgba(11, 21, 32, 0))",
        "radial-highlight": "radial-gradient(circle at 50% 0%, rgba(85, 214, 190, 0.18) 0%, transparent 70%)",
      },
      animation: {
        "pulse-slow": "pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite",
        "float": "float 6s ease-in-out infinite",
        "shimmer": "shimmer 2.5s linear infinite",
      },
      keyframes: {
        float: {
          "0%, 100%": { transform: "translateY(0)" },
          "50%": { transform: "translateY(-6px)" },
        },
        shimmer: {
          "0%": { backgroundPosition: "-200% 0" },
          "100%": { backgroundPosition: "200% 0" },
        }
      }
    },
  },
  plugins: [],
};

export default config;
