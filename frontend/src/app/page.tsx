"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { motion } from "framer-motion";
import { Sparkles, Languages, ArrowRight, BrainCircuit, Rocket, Globe2 } from "lucide-react";

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      router.push("/dashboard");
    }
  }, [router]);

  return (
    <div className="relative min-h-screen overflow-hidden bg-slate-950 flex flex-col items-center justify-center">
      {/* Background Decorative Elements */}
      <div className="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-[120px] animate-pulse" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-purple-500/20 rounded-full blur-[120px] animate-pulse" style={{ animationDelay: '2s' }} />
      </div>

      <main className="relative z-10 container mx-auto px-4 py-20 flex flex-col items-center text-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="flex items-center gap-2 mb-8 px-4 py-1.5 rounded-full bg-white/5 border border-white/10 backdrop-blur-md"
        >
          <Sparkles className="h-4 w-4 text-yellow-500" />
          <span className="text-sm font-medium text-slate-300 tracking-wide uppercase">AI-Powered Language Mastery</span>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          <h1 className="text-6xl md:text-8xl font-black tracking-tight text-white mb-6">
            Linguist<span className="text-primary">AI</span>
          </h1>
          <p className="text-xl md:text-2xl text-slate-400 max-w-2xl mx-auto mb-12 leading-relaxed">
            Experience the future of language learning. Our neural engine builds a personalized path just for you.
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
          className="flex flex-col sm:flex-row gap-4"
        >
          <Button 
            size="lg" 
            onClick={() => router.push("/register")}
            className="h-16 px-10 rounded-full text-lg font-bold shadow-2xl shadow-primary/20 hover:scale-105 transition-transform"
          >
            Start Journey
            <ArrowRight className="ml-2 h-5 w-5" />
          </Button>
          <Button 
            variant="outline" 
            size="lg" 
            onClick={() => router.push("/login")}
            className="h-16 px-10 rounded-full text-lg font-bold border-white/10 bg-white/5 backdrop-blur-md text-white hover:bg-white/10 hover:scale-105 transition-all"
          >
            Sign In
          </Button>
        </motion.div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 1, delay: 0.5 }}
          className="mt-24 grid grid-cols-1 sm:grid-cols-3 gap-8 max-w-4xl w-full"
        >
          <div className="flex flex-col items-center gap-4 p-6 rounded-3xl bg-white/5 border border-white/10 backdrop-blur-sm">
             <div className="h-12 w-12 rounded-2xl bg-primary/20 flex items-center justify-center text-primary">
                <BrainCircuit className="h-6 w-6" />
             </div>
             <div className="text-center">
                <h3 className="text-white font-bold mb-1">Neural Tutor</h3>
                <p className="text-slate-500 text-sm">Adaptive AI that learns your patterns.</p>
             </div>
          </div>
          <div className="flex flex-col items-center gap-4 p-6 rounded-3xl bg-white/5 border border-white/10 backdrop-blur-sm">
             <div className="h-12 w-12 rounded-2xl bg-purple-500/20 flex items-center justify-center text-purple-400">
                <Rocket className="h-6 w-6" />
             </div>
             <div className="text-center">
                <h3 className="text-white font-bold mb-1">Fast Path</h3>
                <p className="text-slate-500 text-sm">Skip what you know, focus on growth.</p>
             </div>
          </div>
          <div className="flex flex-col items-center gap-4 p-6 rounded-3xl bg-white/5 border border-white/10 backdrop-blur-sm">
             <div className="h-12 w-12 rounded-2xl bg-blue-500/20 flex items-center justify-center text-blue-400">
                <Globe2 className="h-6 w-6" />
             </div>
             <div className="text-center">
                <h3 className="text-white font-bold mb-1">10+ Languages</h3>
                <p className="text-slate-500 text-sm">Master global tongues with ease.</p>
             </div>
          </div>
        </motion.div>
      </main>

      {/* Floating Icons for vibe */}
      <motion.div 
        animate={{ 
          y: [0, -20, 0],
          rotate: [0, 5, 0]
        }}
        transition={{ duration: 5, repeat: Infinity, ease: "easeInOut" }}
        className="absolute top-1/4 left-10 text-primary/10 hidden xl:block"
      >
        <Languages size={120} />
      </motion.div>
      <motion.div 
        animate={{ 
          y: [0, 20, 0],
          rotate: [0, -5, 0]
        }}
        transition={{ duration: 6, repeat: Infinity, ease: "easeInOut" }}
        className="absolute bottom-1/4 right-10 text-purple-500/10 hidden xl:block"
      >
        <Globe2 size={120} />
      </motion.div>
    </div>
  );
}

