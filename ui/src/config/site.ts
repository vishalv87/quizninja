export const siteConfig = {
  name: "QuizNinja",
  description: "Test Your Knowledge, Challenge Your Friends",
  url: process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000",
  ogImage: "/og-image.png",
  links: {
    github: "https://github.com/yourusername/quizninja",
  },
  creator: "QuizNinja Team",
};

export type SiteConfig = typeof siteConfig;
