export const humanRelation: Record<string, string> = {
  DEPENDS_ON: "relies on",
  DEPLOYED_AS: "runs as",
  CALLS: "talks to",
  OWNS: "is owned by",
  BELONGS_TO: "belongs to",
  CONTRIBUTED_TO: "helped build",
  MAINTAINS: "looks after",
};

export function humanType(type: string): string {
  switch (type) {
    case "Package": return "Package — a building block like a database or library";
    case "Repository": return "Code area — where the code for one part lives";
    case "Service": return "Running service — does work for users";
    case "Developer": return "Person — writes and maintains code";
    case "Team": return "Team — owns and supports services";
    case "Organization": return "Organization — groups teams together";
    default: return type;
  }
}

export function plainTypeLabel(type: string): string {
  switch (type) {
    case "Package": return "building block";
    case "Repository": return "code area";
    case "Service": return "service";
    case "Developer": return "person";
    case "Team": return "team";
    case "Organization": return "organization";
    default: return type.toLowerCase();
  }
}
