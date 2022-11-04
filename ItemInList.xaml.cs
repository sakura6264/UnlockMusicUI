using System;
using System.Collections.Generic;
using System.Text;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;

namespace UnlockMusicUI
{
    /// <summary>
    /// UserControl1.xaml 的交互逻辑
    /// </summary>
    public partial class ItemInList : UserControl
    {
        private string Textstr;
        public string Description { get; set; }
        public bool IsComplete { get; set; }
        public ItemInList(string Text)
        {
            InitializeComponent();
            this.Text.Text = Text;
            this.Textstr = Text;
            this.Img.Source = (DrawingImage)FindResource("Wait");
            this.Description = string.Empty;
            this.ShowError.IsEnabled = false;
            this.IsComplete = false;
        }
        public void SetState(bool hasError, string description)
        {
            if (hasError)
            {
                this.Description = description;
                this.ShowError.IsEnabled = true;
                this.Img.Source = (DrawingImage)FindResource("Alert");
            }
            else
            {
                this.Img.Source = (DrawingImage)FindResource("OK");
                this.IsComplete = true;
            }
        }

        private void ShowError_Click(object sender, RoutedEventArgs e)
        {
            _ = MessageBox.Show("Error from file:\"" + this.Text.Text + "\"\r\n" + Description,
                "Error", MessageBoxButton.OK, MessageBoxImage.Error);
        }
        public string GetText()
        {
            return this.Textstr;
        }
    }
}
