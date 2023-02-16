using System;
using System.Collections.Generic;
using System.Linq;
using System.IO;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;
using System.ComponentModel;
using Microsoft.Win32;
using Vanara.PInvoke;
using static Vanara.PInvoke.Gdi32;

namespace UnlockMusicUI
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        const string AllSuportExt = "*.666c6163;*.6d3461;*.6d7033;*.6f6767;*.776176;*.aac;*.bkcape;*.bkcflac;*.bkcm4a;*.bkcmp3;*.bkcogg;*.bkcwav;*.bkcwma;*.flac;*.kgm;*.kgma;*.kwm;*.m4a;*.mflac;*.mflac0;*.mflach;*.mgg;*.mgg1;*.mggl;*.mmp4;*.mp3;*.ncm;*.ogg;*.qmc0;*.qmc2;*.qmc3;*.qmc4;*.qmc6;*.qmc8;*.qmcflac;*.qmcogg;*.tkm;*.tm0;*.tm2;*.tm3;*.tm6;*.vpr;*.wav;*.wma;*.x2m;*.x3m;*.xm";
        List<ItemInList> TODOList = new List<ItemInList>();
        public MainWindow()
        {
            InitializeComponent();
            this.TODOs.ItemsSource = TODOList;
            this.TODOs.Items.Refresh();
        }

        private void Adder_Click(object sender, RoutedEventArgs e)
        {
            OpenFileDialog openFileDialog = new OpenFileDialog
            {
                Filter = "Supported Files|" + AllSuportExt + "|All Files|*.*",
                RestoreDirectory = true,
                FilterIndex = 1,
                Multiselect = true
            };
            if (openFileDialog.ShowDialog() == true)
            {
                string[] filename = openFileDialog.FileNames;
                var lst = this.TODOList.Cast<ItemInList>().Select(n => n.GetText());
                foreach (var i in filename)
                {
                    if (lst.Contains(i))
                    {
                        _ = MessageBox.Show("File\"" + i + "\" has already been included in the file list.", "Warning",
                            MessageBoxButton.OK, MessageBoxImage.Warning);
                    }
                    else
                    {
                        this.TODOList.Add(new ItemInList(i));
                    }
                }
                this.TODOs.ItemsSource = TODOList;
                this.TODOs.Items.Refresh();
            }
        }

        private void DirAdder_Click(object sender, RoutedEventArgs e)
        {
            List<string> files = new List<string>();
            void FindFile(DirectoryInfo di)
            {
                FileInfo[] fis = di.GetFiles();
                foreach (var i in fis)
                {
                    files.Add(i.FullName);
                }
                DirectoryInfo[] dis = di.GetDirectories();
                foreach (var i in dis)
                {
                    if (!(i.Name.Equals(".") || i.Name.Equals("..")))
                    {
                        FindFile(i);
                    }
                }
            }
            var strbuf = new Vanara.InteropServices.SafeCoTaskMemString(512);
            Shell32.BROWSEINFO info = new Shell32.BROWSEINFO(HWND.NULL, IntPtr.Zero, "Select Folder",
                Shell32.BrowseInfoFlag.BIF_RETURNONLYFSDIRS | Shell32.BrowseInfoFlag.BIF_USENEWUI | Shell32.BrowseInfoFlag.BIF_DONTGOBELOWDOMAIN | Shell32.BrowseInfoFlag.BIF_NONEWFOLDERBUTTON,
                null, strbuf);
            var res = Shell32.SHBrowseForFolder(info);
            StringBuilder builder = new StringBuilder(512);
            if (!(res.IsNull || res.IsClosed || res.IsInvalid || res.IsEmpty))
            {
                Shell32.SHGetPathFromIDList(res, builder);
                DirectoryInfo dir = new DirectoryInfo(builder.ToString());
                if (dir.Exists)
                {
                    FindFile(dir);
                    var lsts = this.TODOList.Cast<ItemInList>().Select(n => n.GetText());
                    var addlst = files.Where(n => !lsts.Contains(n)).Select(n => new ItemInList(n)).ToList();
                    this.TODOList.AddRange(addlst);
                    this.TODOs.ItemsSource = this.TODOList;
                    this.TODOs.Items.Refresh();
                }
                else
                {
                    _ = MessageBox.Show("File\"" + dir.FullName + "\" has already been included in the file list.", "Warning",
                            MessageBoxButton.OK, MessageBoxImage.Warning);
                }
            }
        }

        private void Remover_Click(object sender, RoutedEventArgs e)
        {
            ItemInList[] lst = this.TODOs.SelectedItems.Cast<ItemInList>().ToArray();
            foreach (var i in lst)
            {
                this.TODOList.Remove(i);
            }
            this.TODOs.ItemsSource = TODOList;
            this.TODOs.Items.Refresh();
        }

        private void Starter_Click(object sender, RoutedEventArgs e)
        {
            if (TODOList.Count == 0)
            {
                _ = MessageBox.Show("Nothing to do.", "Warning",
                           MessageBoxButton.OK, MessageBoxImage.Information);
                return;
            }
            var strbuf = new Vanara.InteropServices.SafeCoTaskMemString(512);
            Shell32.BROWSEINFO info = new Shell32.BROWSEINFO(HWND.NULL, IntPtr.Zero, "Select Folder to Save",
                Shell32.BrowseInfoFlag.BIF_RETURNONLYFSDIRS | Shell32.BrowseInfoFlag.BIF_USENEWUI | Shell32.BrowseInfoFlag.BIF_DONTGOBELOWDOMAIN,
                null, strbuf);
            var res = Shell32.SHBrowseForFolder(info);
            StringBuilder builder = new StringBuilder(512);
            if (!(res.IsNull || res.IsClosed || res.IsInvalid))
            {
                Shell32.SHGetPathFromIDList(res, builder);
                DirectoryInfo dir = new DirectoryInfo(builder.ToString());
                if (dir.Exists)
                {
                    BackgroundWorker worker = new BackgroundWorker
                    {
                        WorkerReportsProgress = true
                    };
                    worker.DoWork += WorkerDoWork;
                    worker.ProgressChanged += WorkerProgressChange;
                    worker.RunWorkerCompleted += WorkerComplete;
                    worker.RunWorkerAsync(dir.FullName);
                }
                else
                {
                    _ = MessageBox.Show("File\"" + dir.FullName + "\" has already been included in the file list.", "Warning",
                            MessageBoxButton.OK, MessageBoxImage.Warning);
                }
            }
        }
        private void Sort_Click(object sender, RoutedEventArgs e)
        {
            this.TODOList = this.TODOList.Cast<ItemInList>().OrderBy(n => n.GetText()).ToList();
            this.TODOs.ItemsSource = TODOList;
            this.TODOs.Items.Refresh();
        }

        private void TODOs_Drop(object sender, DragEventArgs e)
        {
            List<string> files = new List<string>();
            void FindFile(DirectoryInfo di)
            {
                FileInfo[] fis = di.GetFiles();
                foreach (var i in fis)
                {
                    files.Add(i.FullName);
                }
                DirectoryInfo[] dis = di.GetDirectories();
                foreach (var i in dis)
                {
                    if (!(i.Name.Equals(".") || i.Name.Equals("..")))
                    {
                        FindFile(i);
                    }
                }
            }
            foreach (var i in e.Data.GetData(DataFormats.FileDrop) as Array)
            {
                var fi = new FileInfo(i as string);
                if ((fi.Attributes & FileAttributes.Directory) != 0)
                {
                    FindFile(new DirectoryInfo(i as string));
                    var lsts = this.TODOList.Cast<ItemInList>().Select(n => n.GetText());
                    var addlst = files.Where(n => !lsts.Contains(n)).Select(n => new ItemInList(n)).ToList();
                    this.TODOList.AddRange(addlst);
                    files.Clear();
                }
                else
                {
                    string[] lst = this.TODOList.Cast<ItemInList>().Select(n => n.GetText()).ToArray();
                    if (!lst.Contains(i as string))
                    {
                        TODOList.Add(new ItemInList(i as string));
                    }
                }
            }
            this.TODOs.ItemsSource = TODOList;
            this.TODOs.Items.Refresh();
        }

        private void RemoveComp_Click(object sender, RoutedEventArgs e)
        {
            ItemInList[] comp = this.TODOList.Cast<ItemInList>().Where(n => n.IsComplete).ToArray();
            foreach (var i in comp)
            {
                this.TODOList.Remove(i);
            }
            this.TODOs.ItemsSource = TODOList;
            this.TODOs.Items.Refresh();
        }

        private void RemoveAll_Click(object sender, RoutedEventArgs e)
        {
            this.TODOList.Clear();
            this.TODOs.ItemsSource = TODOList;
            this.TODOs.Items.Refresh();
        }
        private void WorkerDoWork(object sender, DoWorkEventArgs e)
        {
            string outpath = e.Argument as string;
            int sum = TODOList.Count;
            int error = 0;
            int pass = 0;
            foreach (var i in TODOList)
            {
                if (i.IsComplete)
                {
                    pass++;
                }
                else
                {
                    int result = Helper.DecFile(Helper.UTF8(i.GetText()), Helper.UTF8(outpath), 0);
                    if (result != 0)
                    {
                        error++;
                    }
                    Helper.ReportArgs reportArgs = new Helper.ReportArgs() { Complete = i, ReturnNum = result };
                    (sender as BackgroundWorker).ReportProgress(0, reportArgs);
                }
            }
            Helper.FinalArgs finalArgs = new Helper.FinalArgs() { Error = error, Pass = pass, Sum = sum };
            e.Result = finalArgs;
        }
        private void WorkerProgressChange(object sender, ProgressChangedEventArgs e)
        {
            var result = e.UserState as Helper.ReportArgs;
            switch (result.ReturnNum)
            {
                case 0:
                    result.Complete.SetState(false, "");
                    break;
                case 1:
                    result.Complete.SetState(true, "No Decoder to Decode This File.");
                    break;
                case -1:
                    result.Complete.SetState(true, "Decode Error.");
                    break;
                case -2:
                    result.Complete.SetState(true, "Write to Path Error.");
                    break;
                default:
                    result.Complete.SetState(true, "Unkown Error.");
                    break;
            }
            this.TODOs.ItemsSource = TODOList;
            this.TODOs.Items.Refresh();

        }
        private void WorkerComplete(object sender, RunWorkerCompletedEventArgs e)
        {
            var result = e.Result as Helper.FinalArgs;
            _ = MessageBox.Show("Decode Completed. Total : " + result.Sum + " Pass : " + result.Pass + " Error : " + result.Error, "Complete",
                            MessageBoxButton.OK, MessageBoxImage.Warning);
        }
    }
}
