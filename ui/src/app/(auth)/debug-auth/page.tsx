"use client";

import { useEffect, useState } from "react";
import { supabase } from "@/lib/supabase/client";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { RefreshCw } from "lucide-react";

export default function DebugAuthPage() {
  const [sessionInfo, setSessionInfo] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [cookies, setCookies] = useState<string>("");

  const checkSession = async () => {
    setLoading(true);
    console.log("=== AUTH DEBUG ===");

    // Check session
    const { data, error } = await supabase.auth.getSession();
    console.log("Session data:", data);
    console.log("Session error:", error);

    // Check user
    const { data: userData, error: userError } = await supabase.auth.getUser();
    console.log("User data:", userData);
    console.log("User error:", userError);

    // Get cookies
    const cookieStr = document.cookie;
    console.log("Cookies:", cookieStr);
    setCookies(cookieStr);

    setSessionInfo({
      session: data.session,
      user: userData.user,
      error: error || userError,
      hasSession: !!data.session,
      hasUser: !!userData.user,
      expiresAt: data.session?.expires_at,
      accessToken: data.session?.access_token ? "***present***" : null,
      refreshToken: data.session?.refresh_token ? "***present***" : null,
    });

    setLoading(false);
  };

  useEffect(() => {
    checkSession();

    // Listen to auth changes
    const { data: authListener } = supabase.auth.onAuthStateChange(
      (event, session) => {
        console.log("Auth state changed:", event, session);
        checkSession();
      }
    );

    return () => {
      authListener.subscription.unsubscribe();
    };
  }, []);

  return (
    <div className="container max-w-4xl py-8">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Auth Debug Page</h1>
          <Button onClick={checkSession} disabled={loading}>
            <RefreshCw className={`h-4 w-4 mr-2 ${loading ? "animate-spin" : ""}`} />
            Refresh
          </Button>
        </div>

        {sessionInfo && (
          <div className="grid gap-4">
            {/* Status Card */}
            <Card>
              <CardHeader>
                <CardTitle>Authentication Status</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="font-medium">Has Session:</span>
                  <Badge variant={sessionInfo.hasSession ? "default" : "destructive"}>
                    {sessionInfo.hasSession ? "Yes ✓" : "No ✗"}
                  </Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="font-medium">Has User:</span>
                  <Badge variant={sessionInfo.hasUser ? "default" : "destructive"}>
                    {sessionInfo.hasUser ? "Yes ✓" : "No ✗"}
                  </Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="font-medium">Access Token:</span>
                  <Badge variant={sessionInfo.accessToken ? "default" : "secondary"}>
                    {sessionInfo.accessToken || "None"}
                  </Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="font-medium">Refresh Token:</span>
                  <Badge variant={sessionInfo.refreshToken ? "default" : "secondary"}>
                    {sessionInfo.refreshToken || "None"}
                  </Badge>
                </div>
              </CardContent>
            </Card>

            {/* User Info */}
            {sessionInfo.user && (
              <Card>
                <CardHeader>
                  <CardTitle>User Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2">
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <span className="font-medium">ID:</span>
                    <span className="font-mono text-xs">{sessionInfo.user.id}</span>
                    <span className="font-medium">Email:</span>
                    <span>{sessionInfo.user.email}</span>
                    <span className="font-medium">Created:</span>
                    <span>{new Date(sessionInfo.user.created_at).toLocaleString()}</span>
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Session Details */}
            {sessionInfo.session && (
              <Card>
                <CardHeader>
                  <CardTitle>Session Details</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2">
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <span className="font-medium">Expires At:</span>
                    <span>
                      {new Date((sessionInfo.expiresAt || 0) * 1000).toLocaleString()}
                    </span>
                    <span className="font-medium">Time Until Expiry:</span>
                    <span>
                      {Math.round(
                        ((sessionInfo.expiresAt || 0) * 1000 - Date.now()) / 1000 / 60
                      )}{" "}
                      minutes
                    </span>
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Cookies */}
            <Card>
              <CardHeader>
                <CardTitle>Browser Cookies</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="bg-muted p-4 rounded text-xs font-mono break-all">
                  {cookies || "No cookies found"}
                </div>
              </CardContent>
            </Card>

            {/* Error Info */}
            {sessionInfo.error && (
              <Card className="border-destructive">
                <CardHeader>
                  <CardTitle className="text-destructive">Error</CardTitle>
                </CardHeader>
                <CardContent>
                  <pre className="text-xs bg-destructive/10 p-4 rounded overflow-auto">
                    {JSON.stringify(sessionInfo.error, null, 2)}
                  </pre>
                </CardContent>
              </Card>
            )}

            {/* Raw Data */}
            <Card>
              <CardHeader>
                <CardTitle>Raw Session Data</CardTitle>
              </CardHeader>
              <CardContent>
                <pre className="text-xs bg-muted p-4 rounded overflow-auto max-h-96">
                  {JSON.stringify(sessionInfo, null, 2)}
                </pre>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
}
