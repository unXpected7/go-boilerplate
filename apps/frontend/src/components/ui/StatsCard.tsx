import { Card, CardContent, CardHeader, CardTitle } from "./card";
import { LucideIcon } from "lucide-react";

interface StatsCardProps {
  title: string;
  value: number | string;
  description?: string;
  icon: LucideIcon;
  trend?: {
    value: number;
    isPositive: boolean;
  };
  className?: string;
}

export function StatsCard({
  title,
  value,
  description,
  icon: Icon,
  trend,
  className = "",
}: StatsCardProps) {
  return (
    <Card className={`${className}`}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {description && (
          <p className="text-xs text-muted-foreground">{description}</p>
        )}
        {trend && (
          <div className={`flex items-center text-xs ${
            trend.isPositive ? "text-green-600" : "text-red-600"
          }`}>
            <span className="mr-1">
              {trend.isPositive ? "↑" : "↓"}
            </span>
            <span>{trend.value}% from last week</span>
          </div>
        )}
      </CardContent>
    </Card>
  );
}