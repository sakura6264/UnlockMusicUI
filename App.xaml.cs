using System;
using System.Collections.Generic;
using System.Configuration;
using System.Data;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Windows;

namespace UnlockMusicUI
{
    /// <summary>
    /// Interaction logic for App.xaml
    /// </summary>
    public partial class App : Application
    {
        MainWindow mainWindow = null;
        private void Application_Startup(object sender, StartupEventArgs e)
        {
            // init dll
            // get exe dir
            var exedir = System.IO.Path.GetDirectoryName(System.Reflection.Assembly.GetExecutingAssembly().Location);
            // convert to byte
            var exedir_utf8 = Helper.UTF8(exedir);
            // init
            Helper.ForceInit(exedir_utf8);
            mainWindow = new MainWindow();
            mainWindow.Show();
        }

        protected override void OnStartup(StartupEventArgs e)
        {
            SplashScreen sp = new SplashScreen("loading.png");
            sp.Show(true);
            base.OnStartup(e);
        }
    }
}
