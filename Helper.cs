using System;
using System.Collections.Generic;
using System.Text;
using System.Runtime.InteropServices;

namespace UnlockMusicUI
{
    public static class Helper
    {
        [DllImport("Um.dll",CharSet = CharSet.Auto,CallingConvention = CallingConvention.Cdecl)]
        public static extern int DecFile(byte[] inFile, byte[] outDir, int SkipNoop);

        [DllImport("Um.dll", CharSet = CharSet.Auto, CallingConvention = CallingConvention.Cdecl)]
        public static extern void ForceInit(byte[] exeDir);
        public static byte[] UTF8(string input)
        {
            var bytes = Encoding.Default.GetBytes(input);
            return Encoding.Convert(Encoding.Default, Encoding.UTF8, bytes);
        }
        public class ReportArgs
        {
            public ItemInList Complete;
            public int ReturnNum;
        }
        public class FinalArgs
        {
            public int Sum;
            public int Pass;
            public int Error;
        }
    }
}
